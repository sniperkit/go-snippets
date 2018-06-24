package rewrite

import (
	"bytes"
	"net/url"
)

type UrlRewriter struct {
	hostRelative bool
	fromHost     string
	pathDepth    int
	to           *url.URL
}

func NewUrlRewriter(from, to string) *UrlRewriter {
	f, err := url.Parse(from)
	if err != nil {
		// TODO - ugh.
		panic(err)
	}

	t, err := url.Parse(to)
	if err != nil {
		// TODO
		panic(err)
	}

	return &UrlRewriter{
		fromHost: f.Host,
		to:       t,
	}
}

// NewRelativeUrlRewriter turns urls that match from's
// hostname into relative urls
func NewRelativeUrlRewriter(from string) *UrlRewriter {
	f, err := url.Parse(from)
	if err != nil {
		// TODO - ugh.
		panic(err)
	}

	return &UrlRewriter{
		fromHost: f.Host,
		to:       &url.URL{},
	}
}

func (urw *UrlRewriter) RewriteString(p string) string {
	return string(urw.Rewrite([]byte(p)))
}

func (urw *UrlRewriter) Rewrite(p []byte) []byte {
	// call to rewrite with empty slice is a no-op
	if len(p) == 0 {
		return nil
	}

	u, err := urw.to.Parse(string(p))
	if err != nil {
		return p
	}

	if u.Host == urw.fromHost {
		u.Host = urw.to.Host
		if u.Scheme != urw.to.Scheme {
			u.Scheme = urw.to.Scheme
		}
	}

	// if we're rewriting to relative urls, ensure
	// empty rewrites to root
	// if urw.to.Host == "" && u.Path == "" {
	// 	u.Path = "/"
	// }

	// relative urls should be "directory relative"
	if u.Host == "" {
		u.Path = "." + u.Path
	}

	if urw.hostRelative {
		u.Scheme = ""
		// rel := u.String()
		return append(urw.pathPrefix(), []byte(u.String())[2:]...)
	}

	return []byte(u.String())
}

func NewHostRelativeUrlRewriter(from string) *UrlRewriter {
	f, err := url.Parse(from)
	if err != nil {
		// TODO - ugh.
		panic(err)
	}

	if f.Path == "" {
		f.Path = "/"
	}

	return &UrlRewriter{
		fromHost:     f.Host,
		hostRelative: true,
		pathDepth:    bytes.Count([]byte(f.Path), []byte{'/'}),
		to:           f,
	}
}

func (urw *UrlRewriter) pathPrefix() []byte {
	return bytes.Repeat([]byte("../"), urw.pathDepth)
}
