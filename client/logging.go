package client

import (
	"log"
	"net/http"
)

// Logging will create a client decorator with logging concerns.
func Logging(l *log.Logger) Decorator {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			l.Println("Client do called: " + r.Host)
			return c.Do(r)
		})
	}
}
