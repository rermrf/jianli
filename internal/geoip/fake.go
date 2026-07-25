package geoip

// fakeResolver is a Resolver that returns canned values regardless of the
// input IP. Lives in a non-test file so other packages' tests can import
// NewFakeResolver to inject fixtures.
type fakeResolver struct {
	country string
	region  string
	city    string
	isp     string
	ok      bool
}

func (f fakeResolver) Resolve(_ string) (country, region, city, isp string, ok bool) {
	return f.country, f.region, f.city, f.isp, f.ok
}

// NewFakeResolver returns a Resolver that always returns the provided
// values. Intended for tests in this and other packages.
func NewFakeResolver(country, region, city, isp string, ok bool) Resolver {
	return fakeResolver{country: country, region: region, city: city, isp: isp, ok: ok}
}
