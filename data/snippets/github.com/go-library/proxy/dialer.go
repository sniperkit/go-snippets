package proxy

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
)

func GetProxyURL() (proxyURL *url.URL, err error) {
	envkeys := []string{
		"HTTP_PROXY", "HTTPS_PROXY",
		"http_proxy", "https_proxy",
	}

	// read environment
	for i := range envkeys {
		env := os.Getenv(envkeys[i])
		if env != "" {
			proxyURL, err = url.Parse(env)
			if err != nil {
				return nil, err
			}
			return proxyURL, nil
		}
	}

	return nil, fmt.Errorf("proxy environment value was not appeared")
}

type Dialer struct {
	ProxyURL *url.URL
}

func New(proxyURL *url.URL) (d *Dialer, err error) {
	d = new(Dialer)
	d.ProxyURL = proxyURL

	return d, nil
}

func (p *Dialer) Dial(network, address string) (conn net.Conn, err error) {

	if network != "tcp" {
		return nil, fmt.Errorf("unsupported network: %s", network)
	}

	if p.ProxyURL == nil {
		return net.Dial(network, address)
	}

	conn, err = net.Dial("tcp", p.ProxyURL.Host)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("CONNECT %s HTTP/1.0\r\n", address)
	if p.ProxyURL.User != nil {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(p.ProxyURL.User.String()))
		authHeader := fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", encodedAuth)
		uri += authHeader
	}
	uri += "\r\n"
	_, err = io.WriteString(conn, uri)
	if err != nil {
		return nil, err
	}

	var (
		headers []string
	)

	headers, err = readHeaders(conn)
	if err != nil {
		return nil, err
	}

	if len(headers) == 0 {
		return nil, fmt.Errorf("proxy header length is zero")
	}

	var (
		ver    string
		status int
	)

	fmt.Sscanf(headers[0], "%s %d", &ver, &status)
	if strings.HasPrefix(ver, "HTTP/") && status >= 200 && status < 300 {
	} else {
		return nil, fmt.Errorf("failed to establish connection")
	}

	return conn, nil
}

func readHeaders(conn net.Conn) ([]string, error) {
	var (
		err     error
		line    []byte
		headers []string
		reader  = bufio.NewReader(conn)
	)

	for {
		line, err = reader.ReadBytes(byte('\n'))
		if err != nil {
			break
		}

		if string(line) == "\r\n" || string(line) == "\n" {
			break
		}

		header := strings.TrimSpace(string(line))
		headers = append(headers, header)
	}

	if err != nil {
		return nil, fmt.Errorf("parse header was failed: %s", err)
	}
	return headers, nil
}
