package proxy

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/websocket"
	"net"
	"strings"
)

func DialWebsocket(url_, origin string) (conn net.Conn, err error) {
	var (
		config *websocket.Config
		d      *Dialer
	)

	config, err = websocket.NewConfig(url_, origin)
	if err != nil {
		return nil, err
	}

	proxyURL, err := GetProxyURL()
	if err != nil {
		return nil, err
	}

	d, err = New(proxyURL)
	if err != nil {
		return nil, err
	}

	addr := config.Location.Host
	if colonPos := strings.LastIndex(config.Location.Host, ":"); colonPos == -1 {
		port := 80
		if config.Location.Scheme == "wss" {
			port = 443
		}
		addr = fmt.Sprintf("%s:%d", config.Location.Host, port)
	}

	conn, err = d.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if config.Location.Scheme == "wss" {
		colonPos := strings.LastIndex(config.Location.Host, ":")
		if colonPos == -1 {
			colonPos = len(config.Location.Host)
		}
		hostname := config.Location.Host[:colonPos]

		config.TlsConfig = &tls.Config{
			ServerName: hostname,
		}
		conn = tls.Client(conn, config.TlsConfig)
	}

	conn, err = websocket.NewClient(config, conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
