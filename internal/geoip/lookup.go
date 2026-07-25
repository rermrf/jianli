// Package geoip resolves IP addresses to geographic regions using a local
// ip2region xdb database. The xdb file is loaded into memory once at
// startup; subsequent lookups are µs-scale and never touch the network.
package geoip

import (
	"errors"
	"net"
	"strings"

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
}

// NewLookup loads the xdb file into memory. Returns a non-nil error if the
// file is missing or corrupt; callers should log and proceed with a nil
// Resolver so the rest of the visitor pipeline still works.
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
	searcher, err := xdb.NewWithBuffer(version, buf)
	if err != nil {
		return nil, err
	}
	return &Lookup{searcher: searcher}, nil
}

// Resolve maps an IPv4/IPv6 string to country/region/city/isp using the
// loaded xdb. Loopback and RFC1918 private addresses short-circuit to
// "本地" / "本地网络" without touching the database.
//
// ip2region xdb v2 returns four (sometimes five) "|"-delimited fields:
//
//	国家 | 省份 | 城市 | ISP | ISO国家码（可选）
//
// We map them to (country, region=省份, city, isp). The optional ISO code
// at parts[4] is ignored. Empty fields surface from ip2region as the
// literal string "0" in older database revisions; cleanZero normalizes
// those to "".
func (l *Lookup) Resolve(ip string) (country, region, city, isp string, ok bool) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", "", "", "", false
	}
	if parsed.IsLoopback() || parsed.IsPrivate() {
		return "本地", "", "本地网络", "", true
	}
	raw, err := l.searcher.Search(ip)
	if err != nil {
		return "", "", "", "", false
	}
	parts := strings.Split(raw, "|")
	if len(parts) < 4 {
		return "", "", "", "", false
	}
	return cleanZero(parts[0]), cleanZero(parts[1]), cleanZero(parts[2]), cleanZero(parts[3]), true
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
