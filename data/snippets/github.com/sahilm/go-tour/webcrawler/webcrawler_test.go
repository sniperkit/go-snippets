package webcrawler_test

import (
	"testing"

	"github.com/sahilm/go-tour/webcrawler"
)

func TestCrawler(t *testing.T) {
	c := webcrawler.NewCrawler("http://golang.org/", 4, fetcher)
	pages, errors := c.Crawl()

	if len(pages) != 4 {
		t.Errorf("got %v pages, want 4", len(pages))
	}

	if len(errors) != 1 {
		t.Errorf("got %v errors, want 1", len(errors))
	}

	errorPage := "http://golang.org/cmd/"

	for _, page := range pages {
		if page.URL == errorPage {
			t.Errorf("got error page: %v in pages", errorPage)
		}
	}

	got := errors[0]
	want := webcrawler.ErrURLNotFound(errorPage)
	if got != want {
		t.Errorf("got error: %v, want error: %v", got, want)
	}
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f *fakeFetcher) Fetch(url string) (*webcrawler.Page, error) {
	if res, ok := (*f)[url]; ok {
		return &webcrawler.Page{URL: url, Body: res.body, Links: res.urls}, nil
	}
	return nil, webcrawler.ErrURLNotFound(url)
}

var fetcher = &fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
