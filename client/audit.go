package client

import (
	"log"
	"net/http"
)

// Audit will create a client decorator with auditing concerns.
func Audit(l *log.Logger) Decorator {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			l.Println("added to audit trail")
			return c.Do(r)
		})
	}
}
