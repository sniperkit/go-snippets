package network

import (
	"regexp"
	"strings"
)

type RouteHandler func(Request) Response
type Route struct {
	Path       string
	Method     string
	Handler    RouteHandler
	AutoAck    bool
	MediaTypes []MediaType
	RegEx      *regexp.Regexp
}

func (r *Route) AutoAcknowledge(ack bool) *Route {
	r.AutoAck = ack

	return r
}

func (r *Route) BindMediaTypes(ms []MediaType) {
	r.MediaTypes = ms
}

func (r *Route) Matches(s string) (bool, map[string]string) {
	matches := r.RegEx.FindAllStringSubmatch(s, -1)
	attrs := make(map[string]string)
	if len(matches) > 0 {
		subExp := r.RegEx.SubexpNames()
		for idx, exp := range subExp {
			attrs[exp] = matches[0][idx]
		}

		return true, attrs
	}

	return false, attrs
}

func MatchingRoute(path string, method string, cf interface{}, routes []*Route) (*Route, map[string]string, error) {
	foundPath := false
	attrs := make(map[string]string)
	for _, route := range routes {
		match, att := route.Matches(path)
		if match {
			attrs = att
			foundPath = true
			if route.Method == method {
				if len(route.MediaTypes) > 0 {

					if cf == nil {
						return route, attrs, ERR_UNSUPPORTED_CONTENT_FORMAT
					}

					foundMediaType := false
					for _, o := range route.MediaTypes {
						if uint32(o) == cf {
							foundMediaType = true
							break
						}
					}

					if !foundMediaType {
						return route, attrs, ERR_UNSUPPORTED_CONTENT_FORMAT
					}
				}
				return route, attrs, nil
			}
		}
	}

	if foundPath {
		return &Route{}, attrs, ERR_NO_MATCHING_METHOD
	} else {
		return &Route{}, attrs, ERR_NO_MATCHING_ROUTE
	}
}

func CreateCompilableRoutePath(pre string) (*regexp.Regexp, error) {
	re, _ := regexp.Compile(`{[a-z]+}`)

	matches := re.FindAllStringSubmatch(pre, -1)

	pre = "^" + pre
	for _, b := range matches {
		origAttr := b[0]
		attr := strings.Replace(strings.Replace(origAttr, "{", "", -1), "}", "", -1)
		frag := `(?P<` + attr + `>\w+)`
		pre = strings.Replace(pre, origAttr, frag, -1)
	}
	pre += "$"
	re, err := regexp.Compile(pre)

	return re, err
}

func CreateNewRoute(path string, method string, fn RouteHandler) *Route {
	re, _ := CreateCompilableRoutePath(path)

	return &Route{
		AutoAck: false,
		Path:    path,
		Method:  method,
		Handler: fn,
		RegEx:   re,
	}

}
