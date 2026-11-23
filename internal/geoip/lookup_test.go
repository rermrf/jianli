package geoip

import (
	"encoding/binary"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

func TestNewLookupReturnsErrorForEmptyPath(t *testing.T) {
	_, err := NewLookup("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

// Write a real xdb containing a single IP so tests exercise the library's
// binary search, address-family checks, and database-format parsing together.
func writeTestXDB(t *testing.T, path, ip, raw string, version *xdb.Version, legacy bool) {
	t.Helper()
	indexPtr := xdb.HeaderInfoLength + xdb.VectorIndexRows*xdb.VectorIndexCols*xdb.VectorIndexSize
	dataPtr := indexPtr + version.SegmentIndexSize
	buf := make([]byte, dataPtr+len(raw))
	structure := uint16(xdb.Structure30)
	if legacy {
		structure = xdb.Structure20
	}
	binary.LittleEndian.PutUint16(buf, structure)
	binary.LittleEndian.PutUint16(buf[2:], uint16(xdb.VectorIndexPolicy))
	binary.LittleEndian.PutUint32(buf[8:], uint32(indexPtr))
	binary.LittleEndian.PutUint32(buf[12:], uint32(indexPtr))
	binary.LittleEndian.PutUint16(buf[16:], uint16(version.Id))
	binary.LittleEndian.PutUint16(buf[18:], 4)
	ipBytes := net.ParseIP(ip).To16()
	if version == xdb.IPv4 {
		ipBytes = net.ParseIP(ip).To4()
	}
	vectorPtr := xdb.HeaderInfoLength + (int(ipBytes[0])*xdb.VectorIndexCols+int(ipBytes[1]))*xdb.VectorIndexSize
	binary.LittleEndian.PutUint32(buf[vectorPtr:], uint32(indexPtr))
	binary.LittleEndian.PutUint32(buf[vectorPtr+4:], uint32(indexPtr))
	for offset := indexPtr; offset < indexPtr+2*version.Bytes; offset += version.Bytes {
		if version == xdb.IPv4 {
			binary.LittleEndian.PutUint32(buf[offset:], binary.BigEndian.Uint32(ipBytes))
		} else {
			copy(buf[offset:], ipBytes)
		}
	}
	binary.LittleEndian.PutUint16(buf[indexPtr+2*version.Bytes:], uint16(len(raw)))
	binary.LittleEndian.PutUint32(buf[indexPtr+2*version.Bytes+2:], uint32(dataPtr))
	copy(buf[dataPtr:], raw)
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLookupResolvesDatabaseFormats(t *testing.T) {
	for _, tc := range []struct {
		name, ip, query, raw string
		version              *xdb.Version
		legacy               bool
	}{
		{"current IPv4 with ISO", "112.17.45.201", "112.17.45.201", "中国|浙江|杭州|移动|CN", xdb.IPv4, false},
		{"current IPv4 without ISO", "112.17.45.201", "112.17.45.201", "中国|浙江|杭州|移动", xdb.IPv4, false},
		{"legacy IPv4", "112.17.45.201", "112.17.45.201", "中国|0|浙江|杭州|移动", xdb.IPv4, true},
		{"IPv4 mapped IPv6", "112.17.45.201", "::ffff:112.17.45.201", "中国|浙江|杭州|移动|CN", xdb.IPv4, false},
		{"current IPv6", "240e:3b7:3272:d8d0:db09:c067:8d59:539e", "240e:3b7:3272:d8d0:db09:c067:8d59:539e", "中国|浙江|杭州|移动|CN", xdb.IPv6, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.xdb")
			writeTestXDB(t, path, tc.ip, tc.raw, tc.version, tc.legacy)
			lookup, err := NewLookup(path)
			if err != nil {
				t.Fatal(err)
			}
			country, region, city, isp, ok := lookup.Resolve(tc.query)
			if !ok || country != "中国" || region != "浙江" || city != "杭州" || isp != "移动" {
				t.Fatalf("Resolve(%q) = %q, %q, %q, %q, %v", tc.query, country, region, city, isp, ok)
			}
		})
	}
}

func TestNewResolverUsesBundledDatabasesAndBothFamilies(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.Mkdir("geoip", 0o755); err != nil {
		t.Fatal(err)
	}
	v4, v6 := "112.17.45.201", "240e:3b7:3272:d8d0:db09:c067:8d59:539e"
	writeTestXDB(t, filepath.Join("geoip", "ip2region_v4.xdb"), v4, "中国|浙江|杭州|移动|CN", xdb.IPv4, false)
	writeTestXDB(t, filepath.Join("geoip", "ip2region_v6.xdb"), v6, "中国|广东|深圳|电信|CN", xdb.IPv6, false)
	r, err := NewResolver("./data/ip2region_v4.xdb", "./data/ip2region_v6.xdb")
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				for _, tc := range []struct{ ip, city string }{{v4, "杭州"}, {v6, "深圳"}, {"127.0.0.1", "本地网络"}} {
					_, _, city, _, ok := r.Resolve(tc.ip)
					if !ok || city != tc.city {
						t.Errorf("Resolve(%q) city=%q ok=%v", tc.ip, city, ok)
					}
				}
			}
		}()
	}
	wg.Wait()
}

func TestNewResolverRetainsIPv4WhenIPv6DatabaseIsMissing(t *testing.T) {
	t.Chdir(t.TempDir())
	writeTestXDB(t, "v4.xdb", "112.17.45.201", "中国|浙江|杭州|移动|CN", xdb.IPv4, false)
	r, err := NewResolver("v4.xdb", "missing-v6.xdb")
	if err == nil {
		t.Fatal("expected an error for the missing IPv6 database")
	}
	if _, _, city, _, ok := r.Resolve("112.17.45.201"); !ok || city != "杭州" {
		t.Fatalf("IPv4 lookup was lost: city=%q ok=%v", city, ok)
	}
	if _, _, _, _, ok := r.Resolve("240e:3b7::1"); ok {
		t.Fatal("expected unresolved IPv6 without its database")
	}
}

func TestLookupDoesNotResolveEmptyGeography(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.xdb")
	writeTestXDB(t, path, "8.8.8.8", "0|0|0|0|US", xdb.IPv4, false)
	r, err := NewLookup(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, _, ok := r.Resolve("8.8.8.8"); ok {
		t.Fatal("empty geography must remain unresolved")
	}
}

func TestNewLookupReturnsErrorForMissingFile(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "does-not-exist.xdb")
	_, err := NewLookup(missingPath)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestNewLookupRejectsTruncatedDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "truncated.xdb")
	writeTestXDB(t, path, "8.8.8.8", "美国|加利福尼亚|山景城|Google|US", xdb.IPv4, false)
	// Keep the full header and vector index but truncate the segment index.
	if err := os.Truncate(path, int64(xdb.HeaderInfoLength+xdb.VectorIndexRows*xdb.VectorIndexCols*xdb.VectorIndexSize)); err != nil {
		t.Fatal(err)
	}
	if _, err := NewLookup(path); err == nil {
		t.Fatal("expected an error before an incomplete database can be queried")
	}
}

func TestCleanZeroNormalizesPlaceholderToEmpty(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"0", ""},
		{" 0 ", ""},
		{"", ""},
		{"中国", "中国"},
		{" 浙江 ", "浙江"},
	}
	for _, tc := range cases {
		if got := cleanZero(tc.in); got != tc.want {
			t.Fatalf("cleanZero(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestResolveLoopbackAndPrivateDoNotTouchSearcher(t *testing.T) {
	l := &Lookup{searcher: nil}

	cases := []struct {
		name        string
		ip          string
		wantOK      bool
		wantCity    string
		wantCountry string
	}{
		{"loopback ipv4", "127.0.0.1", true, "本地网络", "本地"},
		{"private 10/8", "10.0.0.1", true, "本地网络", "本地"},
		{"private 192.168/16", "192.168.1.1", true, "本地网络", "本地"},
		{"loopback ipv6", "::1", true, "本地网络", "本地"},
		{"malformed", "not-an-ip", false, "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			country, _, city, _, ok := l.Resolve(tc.ip)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if city != tc.wantCity {
				t.Fatalf("city = %q, want %q", city, tc.wantCity)
			}
			if country != tc.wantCountry {
				t.Fatalf("country = %q, want %q", country, tc.wantCountry)
			}
		})
	}
}

func TestLoopbackResolverDoesNotResolvePublicIP(t *testing.T) {
	r := NewLoopbackResolver()
	_, _, _, _, ok := r.Resolve("8.8.8.8")
	if ok {
		t.Fatal("loopback-only resolver should report ok=false for public IP")
	}
}

func TestFakeResolverImplementsResolver(t *testing.T) {
	var r Resolver = NewFakeResolver("中国", "浙江", "杭州", "电信", true)
	country, region, city, isp, ok := r.Resolve("114.114.114.114")
	if !ok {
		t.Fatal("expected ok=true from fake")
	}
	if country != "中国" || region != "浙江" || city != "杭州" || isp != "电信" {
		t.Fatalf("unexpected fake fields: %s %s %s %s", country, region, city, isp)
	}
}
