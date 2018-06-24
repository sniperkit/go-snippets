package client

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func BuildUserAgent(progname, version string) string {
	return fmt.Sprintf("%s/%s", progname, version)
}

func buildUrl(c *Config) string {
	//return fmt.Sprintf("%s://%s:%s%s?%s", c.Scheme, c.Host, strconv.Itoa(c.Port), c.Path, query.Encode())
	return fmt.Sprintf("%s://%s:%s%s", c.Scheme, c.Host, strconv.Itoa(c.Port), c.Path)
}

type URL struct {
	Scheme   string
	Opaque   string
	User     *url.Userinfo
	Host     string
	Port     int
	Path     string
	RawPath  string
	RawQuery string
	Fragment string
}

func ParseUrl(rawurl string) (*URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return &URL{}, err
	}

	host := u.Host
	port := 80

	if strings.Contains(u.Host, ":") {
		r := strings.Split(u.Host, ":")

		host = r[0]
		port, _ = strconv.Atoi(r[1])
	} else {
		if u.Scheme == "https" {
			port = 443
		}
	}

	return &URL{
		Scheme:   u.Scheme,
		Opaque:   u.Opaque,
		User:     u.User,
		Host:     host,
		Port:     port,
		Path:     u.Path,
		RawPath:  u.RawPath,
		RawQuery: u.RawQuery,
		Fragment: u.Fragment,
	}, nil
}
