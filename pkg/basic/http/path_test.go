package http

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

func TestPathJoin(t *testing.T) {
	u, err := url.Parse("../../..//search?q=dotnet")
	if err != nil {
		log.Fatal(err)
	}
	base, err := url.Parse("http://example.com/directory/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(base.ResolveReference(u))

}
