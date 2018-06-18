package client

import (
	"math/rand"
	"net/http"
)

type Director func(*http.Request)

func LoadBalancing(dir Director) Decorator {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			dir(r)
			return c.Do(r)
		})
	}
}

func RoundRobin(robin uint64, backends ...string) Director {
	return func(r *http.Request) {
		if len(backends) > 0 {
			r.URL.Host = backends[rand.Intn(len(backends))]
		}
	}
}
