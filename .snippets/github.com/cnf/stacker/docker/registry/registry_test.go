package registry

import (
	"testing"
)

func TestNameFromServerAddress(t *testing.T) {
	urls := map[string]string{
		"192.168.0.1:5000": "192.168.0.1:5000",
		"http://www.example.com:5000/": "www.example.com:5000",
		"https://example.com:443/": "example.com:443",
		"http://example.com:5000/path": "example.com:5000",
		"http://example.com/path": "example.com",
		"lala://example.com:123/": "",
	}
	for k, v := range urls {
		result, err := nameFromServerAddress(k)
		if v == "" && err != ErrNotARegistryAddress {
			t.Errorf("%s is not `ErrNotARegistryAddress`")
		} else if result != v {
			t.Errorf("%s is not %s", result, v)
		}
	}
}
