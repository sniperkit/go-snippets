package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

// middleware handles request logging
func middleware(handler http.HandlerFunc) http.HandlerFunc {
	// no-auth middware func
	return func(w http.ResponseWriter, r *http.Request) {
		// poor man's logging:
		fmt.Println(r.Method, r.URL.Path, time.Now())

		// If this server is operating behind a proxy, but we still want to force
		// users to use https, cfg.ProxyForceHttps == true will listen for the common
		// X-Forward-Proto & redirect to https
		if cfg.ProxyForceHttps {
			if r.Header.Get("X-Forwarded-Proto") == "http" {
				w.Header().Set("Connection", "close")
				url := "https://" + r.Host + r.URL.String()
				http.Redirect(w, r, url, http.StatusMovedPermanently)
				return
			}
		}

		// TODO - Strict Transport config?
		// if cfg.TLS {
		// 	// If TLS is enabled, set 1 week strict TLS, 1 week for now to prevent catastrophic mess-ups
		// 	w.Header().Add("Strict-Transport-Security", "max-age=604800")
		// }

		addCORSHeaders(w, r)

		// remove domainless cookie if one is found
		if cfg.UserCookieDomain != "" {
			sess, err := sessions.NewCookieStore([]byte(cfg.SessionSecret)).Get(r, cfg.UserCookieKey)
			if err != nil {
				sess.Options.MaxAge = -1
				sess.Save(r, w)
			}
		}

		ctx := r.Context()
		u, _ := userFromRequest(appDB, r)
		if u != nil {
			ctx = context.WithValue(ctx, "user", u)
		}
		handler(w, r.WithContext(ctx))
	}
}

// authMiddleware adds http basic auth if configured
func authMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	// return auth middleware if configuration settings are present
	if cfg.HttpAuthUsername != "" && cfg.HttpAuthPassword != "" {
		return func(w http.ResponseWriter, r *http.Request) {
			// poor man's logging:
			fmt.Println(r.Method, r.URL.Path, time.Now())

			// If this server is operating behind a proxy, but we still want to force
			// users to use https, cfg.ProxyForceHttps == true will listen for the common
			// X-Forward-Proto & redirect to https
			if cfg.ProxyForceHttps {
				if r.Header.Get("X-Forwarded-Proto") == "http" {
					w.Header().Set("Connection", "close")
					url := "https://" + r.Host + r.URL.String()
					http.Redirect(w, r, url, http.StatusMovedPermanently)
					return
				}
			}

			user, pass, ok := r.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(cfg.HttpAuthUsername)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(cfg.HttpAuthPassword)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// TODO - Strict Transport config?
			// if cfg.TLS {
			// 	// If TLS is enabled, set 1 week strict TLS, 1 week for now to prevent catastrophic mess-ups
			// 	w.Header().Add("Strict-Transport-Security", "max-age=604800")
			// }
			addCORSHeaders(w, r)

			ctx := r.Context()
			u, _ := userFromRequest(appDB, r)
			if u != nil {
				ctx = context.WithValue(ctx, "user", u)
			}
			handler(w, r.WithContext(ctx))
		}
	}

	// no-auth middware func
	return middleware(handler)
}

// addCORSHeaders adds CORS header info for whitelisted servers
func addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, o := range cfg.AllowedOrigins {
		if origin == o {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			return
		}
	}
}
