// Package geoip resolves IP addresses to geographic regions using a local
// ip2region xdb database. The xdb file is loaded into memory once at
// startup; subsequent lookups are µs-scale and never touch the network.
package geoip

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// Resolver is the abstraction handlers depend on. The production
// implementation is *Lookup; tests inject a fake.
type Resolver interface {
	// Resolve returns country/region/city/isp for the given IP. ok=false
	// signals "could not parse / not found"; callers should leave the
	// fields empty in storage so the UI can fall back to "未知".
	Resolve(ip string) (country, region, city, isp string, ok bool)
}

// Lookup is the ip2region-backed implementation of Resolver.
type Lookup struct {
	searcher *xdb.Searcher
	legacy   bool
	mu       sync.Mutex // xdb.Searcher mutates its IO counter even with a memory buffer.
}

// NewLookup loads the xdb file into memory. Returns a non-nil error if the
// file is missing or corrupt; callers can use NewResolver for fallback.
func NewLookup(xdbPath string) (*Lookup, error) {
	if xdbPath == "" {
		return nil, errors.New("geoip: empty xdb path")
	}
	header, err := xdb.LoadHeaderFromFile(xdbPath)
	if err != nil {
		return nil, err
	}
	version, err := xdb.VersionFromHeader(header)
	if err != nil {
		return nil, err
	}
	buf, err := xdb.LoadContentFromFile(xdbPath)
	if err != nil {
		return nil, err
	}
	minIndexPtr := uint32(xdb.HeaderInfoLength + xdb.VectorIndexRows*xdb.VectorIndexCols*xdb.VectorIndexSize)
	if len(buf) < int(minIndexPtr) || header.StartIndexPtr < minIndexPtr || header.StartIndexPtr > header.EndIndexPtr ||
		uint64(header.EndIndexPtr)+uint64(version.SegmentIndexSize) > uint64(len(buf)) {
		return nil, fmt.Errorf("geoip: truncated xdb file %s", xdbPath)
	}
	searcher, err := xdb.NewWithBuffer(version, buf)
	if err != nil {
		return nil, err
	}
	return &Lookup{searcher: searcher, legacy: header.Version == xdb.Structure20}, nil
}

type dualResolver struct {
	lookups map[int]*Lookup
}

// NewResolver loads both address families. Docker bundles databases in
// ./geoip, outside the mounted ./data directory; configured files take priority.
// A usable resolver is always returned, with errors for unavailable databases.
func NewResolver(paths ...string) (Resolver, error) {
	r := &dualResolver{lookups: make(map[int]*Lookup)}
	var errs []error
	for _, path := range paths {
		if path == "" {
			continue
		}
		lookup, err := NewLookup(path)
		if err != nil {
			bundled := filepath.Join("geoip", filepath.Base(path))
			lookup, err = NewLookup(bundled)
			if err != nil {
				errs = append(errs, fmt.Errorf("geoip: load %s (fallback %s): %w", path, bundled, err))
				continue
			}
		}
		r.lookups[lookup.searcher.IPVersion().Id] = lookup
	}
	return r, errors.Join(errs...)
}

func (r *dualResolver) Resolve(ip string) (country, region, city, isp string, ok bool) {
	if country, region, city, isp, ok = (loopbackResolver{}).Resolve(ip); ok {
		return
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", "", "", "", false
	}
	version := xdb.IPv6VersionNo
	if parsed.To4() != nil {
		version = xdb.IPv4VersionNo
	}
	if lookup := r.lookups[version]; lookup != nil {
		return lookup.Resolve(parsed.String())
	}
	return "", "", "", "", false
}

// Resolve maps an IPv4/IPv6 string to country/region/city/isp using the
// loaded xdb. Loopback and RFC1918 private addresses short-circuit to
// "本地" / "本地网络" without touching the database.
//
// Current databases return four fields plus an optional ISO country code:
//
//	国家 | 省份 | 城市 | ISP | ISO国家码（可选）
//
// Older structure-2 databases return 国家|区域|省份|城市|ISP instead.
// The header distinguishes these formats; field count alone is ambiguous.
func (l *Lookup) Resolve(ip string) (country, region, city, isp string, ok bool) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", "", "", "", false
	}
	if parsed.IsLoopback() || parsed.IsPrivate() {
		return "本地", "", "本地网络", "", true
	}
	if l.searcher == nil {
		return "", "", "", "", false
	}
	l.mu.Lock()
	raw, err := l.searcher.Search(parsed.String())
	l.mu.Unlock()
	if err != nil {
		return "", "", "", "", false
	}
	parts := strings.Split(raw, "|")
	if l.legacy {
		if len(parts) < 5 {
			return "", "", "", "", false
		}
		parts = []string{parts[0], parts[2], parts[3], parts[4]}
	}
	if len(parts) < 4 {
		return "", "", "", "", false
	}
	country, region, city, isp = cleanZero(parts[0]), cleanZero(parts[1]), cleanZero(parts[2]), cleanZero(parts[3])
	return country, region, city, isp, country != "" || region != "" || city != ""
}

func cleanZero(s string) string {
	s = strings.TrimSpace(s)
	if s == "0" {
		return ""
	}
	return s
}

// loopbackResolver is a fallback Resolver used when the xdb file is missing
// or unreadable. It still recognizes loopback / private IPs (so visits from
// localhost or LAN show "本地网络" rather than "未知") but returns "not
// found" for any public IP. main.go falls back to this when NewLookup
// fails so handlers can always rely on a non-nil Resolver.
type loopbackResolver struct{}

func (loopbackResolver) Resolve(ip string) (country, region, city, isp string, ok bool) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", "", "", "", false
	}
	if parsed.IsLoopback() || parsed.IsPrivate() {
		return "本地", "", "本地网络", "", true
	}
	return "", "", "", "", false
}

// NewLoopbackResolver returns a Resolver that only recognizes loopback and
// private addresses. Used as a fallback when the ip2region xdb is missing.
func NewLoopbackResolver() Resolver {
	return loopbackResolver{}
}
