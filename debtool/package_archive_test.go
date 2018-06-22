package debtool_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/strothj/alpine-debtool/debtool"
)

func TestPackageArchive(t *testing.T) {
	t.Skip()
	pa := debtool.PackageArchive{
		Type:         debtool.BinaryArchive,
		Root:         &url.URL{Path: "http://ftp.debian.org/debian"},
		Distribution: "squeeze",
		Components:   []string{"main", "contrib", "non-free"},
	}
	t.Log(pa.Root.String())
	resp, err := http.Get(pa.Root.Path)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
