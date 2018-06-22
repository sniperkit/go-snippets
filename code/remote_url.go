package code

import (
	"regexp"
)

// https://regex101.com/r/LEikJj/2
var gitRemoteUrlRegexp = regexp.MustCompile(`^(?:(?:(?P<protocol>ssh|https?)://)?)(?:(?P<username>\S+)@)?(?P<hostname>(?:[^/:])+)(?P<path_separator>/|:)(?P<path>\S+?)(?P<extension>(?:\.git)?)$`)

// TODO: add password
type remoteURL struct {
	protocol      string
	username      string
	hostname      string
	pathSeparator string
	path          string
	extension     string
}

func (r *remoteURL) String() string {
	var res string
	if r.protocol != "" {
		res += r.protocol + "://"
	}
	if r.username != "" {
		res += r.username + "@"
	}
	res += r.hostname + r.pathSeparator + r.path
	if r.extension != "" {
		res += r.extension
	}

	return res
}

func parseRemoteURL(url string) *remoteURL {
	match := gitRemoteUrlRegexp.FindStringSubmatch(url)
	res := make(map[string]string)
	for i, name := range gitRemoteUrlRegexp.SubexpNames() {
		if i != 0 {
			res[name] = match[i]
		}
	}
	return &remoteURL{
		protocol:      res["protocol"],
		username:      res["username"],
		hostname:      res["hostname"],
		pathSeparator: res["path_separator"],
		path:          res["path"],
		extension:     res["extension"],
	}
}
