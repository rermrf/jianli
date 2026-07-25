package geoip

import (
	"path/filepath"
	"testing"
)

func TestNewLookupReturnsErrorForEmptyPath(t *testing.T) {
	_, err := NewLookup("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestNewLookupReturnsErrorForMissingFile(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "does-not-exist.xdb")
	_, err := NewLookup(missingPath)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
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
