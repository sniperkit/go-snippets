// Resources is a package for extracting urls of dependant files for displaying a web page.
// It's primary use is for establishing the list of assets an archive would need to cache in order
// to properly display a target resource
package resources

import (
	"bytes"

	"github.com/datatogether/warc"
	"golang.org/x/net/html"
)

type Extractor struct {
	Tags map[string][]string
}

func NewExtractor() *Extractor {
	return &Extractor{
		Tags: extractTags(),
	}
}

func (e Extractor) ExtractResponseUrls(rec *warc.Record) ([]string, error) {
	rt := rec.Headers[warc.FieldNameWARCIdentifiedPayloadType]
	switch rt {
	case "text/html; charset=utf-8":
		return e.ExtractHtmlResources(rec)
	default:
		// TODO - for now we just return nothing.
		return nil, nil
		// return nil, fmt.Errorf("can't extract urls from record with content type: '%s'", rt)
	}
	return nil, nil
}

func (e Extractor) ExtractHtmlResources(rec *warc.Record) (urls []string, err error) {
	var p []byte
	p, err = rec.Body()
	if err != nil {
		return
	}
	urls = []string{}
	rdr := bytes.NewReader(p)
	tokenizer := html.NewTokenizer(rdr)

	for {
		tt := tokenizer.Next()
		// token := tokenizer.Token()
		switch tt {
		// case html.TextToken:
		// case html.CommentToken:
		case html.ErrorToken:
			// ErrorToken means that an error occurred during tokenization.
			// most common is end-of-file (EOF)
			if tokenizer.Err().Error() == "EOF" {
				return urls, nil
			}
			return urls, tokenizer.Err()
		case html.StartTagToken:
			name, hasAttr := tokenizer.TagName()
			token := html.Token{
				Type: html.StartTagToken,
				Data: string(name),
			}
			if hasAttr {
				e.extractUrls(&token, tokenizer, &urls)
			}
			continue
		case html.SelfClosingTagToken:
			name, hasAttr := tokenizer.TagName()
			token := html.Token{
				Type: html.SelfClosingTagToken,
				Data: string(name),
			}
			if hasAttr {
				e.extractUrls(&token, tokenizer, &urls)
			}
			continue
		}

	}

	return urls, nil
}

func (e Extractor) extractUrls(t *html.Token, tok *html.Tokenizer, urls *[]string) {
	attrs := e.Tags[t.Data]
	for {
		key, val, more := tok.TagAttr()
		lk := string(bytes.ToLower(key))
		for _, t := range attrs {
			if t == lk {
				*urls = append(*urls, string(val))
			}
		}
		if !more {
			return
		}
	}
	return
}

func extractTags() map[string][]string {
	// oe := PrefixRewriter{Prefix: []byte("oe_")}
	// im := PrefixRewriter{Prefix: []byte("im_")}
	// if_ := PrefixRewriter{Prefix: []byte("if_")}
	// fr_ := PrefixRewriter{Prefix: []byte("fr_")}
	// js_ := PrefixRewriter{Prefix: []byte("js_")}

	return map[string][]string{
		// "a":          []string{"href"},
		"applet":     []string{"codebase", "archive"},
		"area":       []string{"href"},
		"audio":      []string{"src"},
		"base":       []string{"href"},
		"blockquote": []string{"cite"},
		"body":       []string{"background"},
		// "button":     []string{"formaction"},
		"command": []string{"icon"},
		"del":     []string{"cite"},
		"embed":   []string{"src"},
		// "head":      []string {"": defmod}, // for heang
		// "iframe": []string{"src"},
		"image": []string{"src", "xlink:href"},
		"img":   []string{"src", "srcset"},
		"ins":   []string{"cite"},
		// "input":  []string{"src", "formaction"},
		"input": []string{"src"},
		"form":  []string{"action"},
		"frame": []string{"src"},
		"link":  []string{"href"},
		// "meta":   []string{"content"},
		"object": []string{"codebase", "data"},
		"param":  []string{"value"},
		"q":      []string{"cite"},
		"ref":    []string{"href"},
		"script": []string{"src"},
		"source": []string{"src"},
		"video":  []string{"src", "poster"},
	}
}
