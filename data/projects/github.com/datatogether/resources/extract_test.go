package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/datatogether/warc"
)

func TestExtractResponseUrls(t *testing.T) {
	recs, err := readTestWarc("test_extract_response_urls.warc")
	if err != nil {
		t.Error(err.Error())
		return
	}

	e := NewExtractor()

	cases := []struct {
		rec  *warc.Record
		urls []string
		err  string
	}{
		{recs[1], []string{
			"http://datatogether.org/",
			"./css/style.css",
			"./img/favicon.ico",
			"https://s3.amazonaws.com/datatogether/svg/lines_left.svg",
			"https://s3.amazonaws.com/datatogether/svg/lines_right.svg",
			"https://s3.us-east-2.amazonaws.com/static.archivers.space/add-metadata.png",
			"https://s3.amazonaws.com/datatogether/svg/lines_left.svg",
			"https://s3.amazonaws.com/datatogether/svg/lines_right.svg",
			"./js/site.js",
		}, ""},
	}

	for i, c := range cases {
		urls, err := e.ExtractResponseUrls(c.rec)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d marshal error mismatch: expected: %s, got: %s", i, c.err, err)
			continue
		}

		if len(urls) != len(c.urls) {
			for j, u := range urls {
				fmt.Printf("case %d returned urls:", i)
				fmt.Printf("%d: %s\n", j, u)
			}
			t.Errorf("case %d url length mismatch. expected: %d, got: %d", i, len(c.urls), len(urls))
			continue
		}

		for j, u := range urls {
			if c.urls[j] != u {
				t.Errorf("case %d url %d mistmatch. expected: %s, got: %s", i, j, c.urls[j], u)
				break
			}
		}

	}
}

func readTestWarc(file string) (warc.Records, error) {
	f, err := os.Open(filepath.Join("testdata", file))
	if err != nil {
		return nil, err
	}
	r, err := warc.NewReader(f)
	if err != nil {
		return nil, err
	}
	return r.ReadAll()
}
