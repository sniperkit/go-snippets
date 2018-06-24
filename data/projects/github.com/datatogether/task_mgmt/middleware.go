package main

import (
	"crypto/tls"
	"net/http"
	"time"
)

func init() {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}
}

// middleware handles request logging
func middleware(handler http.HandlerFunc) http.HandlerFunc {
	// no-auth middware func
	return func(w http.ResponseWriter, r *http.Request) {
		// poor man's logging:
		log.Infoln(r.Method, r.URL.Path, time.Now())

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

		addCORSHeaders(w, r)

		// TODO - Strict Transport config?
		// if cfg.TLS {
		// 	// If TLS is enabled, set 1 week strict TLS, 1 week for now to prevent catastrophic mess-ups
		// 	w.Header().Add("Strict-Transport-Security", "max-age=604800")
		// }
		handler(w, r)
	}
}

// authMiddleware checks for github auth
// TODO - this is a carry-over from a former implementation of task_mgmt
// that was specific to executing the kiwix zim task it should be shifted
// over to some sort of permissions service

// func authMiddleware(handler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		token := r.FormValue("access_token")
// 		c, err := r.Cookie(cfg.UserCookieKey)
// 		if err != nil {
// 			log.Info(err.Error())
// 		}

// 		// we gots no login info, so login required
// 		if c == nil && token == "" {
// 			msg := fmt.Sprintf("github login required: %s/oauth/github?redirect=%s", cfg.IdentityServerUrl, cfg.UrlRoot)
// 			apiutil.WriteMessageResponse(w, msg, nil)
// 			return
// 		}

// 		req, err := http.NewRequest("GET", fmt.Sprintf("%s/session/oauth/github/repoaccess?access_token=%s&owner=%s&repo=%s", cfg.IdentityServerUrl, token, cfg.GithubRepoOwner, cfg.GithubRepoName), nil)
// 		if err != nil {
// 			log.Info(err.Error())
// 			apiutil.WriteErrResponse(w, http.StatusInternalServerError, fmt.Errorf("error contacting identity server: %s", err.Error()))
// 			return
// 		}

// 		req.AddCookie(c)
// 		res, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			log.Info(err.Error())
// 			apiutil.WriteErrResponse(w, http.StatusInternalServerError, fmt.Errorf("error contacting identity server: %s", err.Error()))
// 			return
// 		}
// 		defer res.Body.Close()

// 		data := map[string]interface{}{}
// 		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
// 			log.Info(err.Error())
// 			apiutil.WriteErrResponse(w, http.StatusInternalServerError, fmt.Errorf("error contacting identity server: %s", err.Error()))
// 			return
// 		}

// 		// User Needs github added to their account
// 		if res.StatusCode == http.StatusUnauthorized {
// 			msg := fmt.Sprintf("github login required: %s/oauth/github?redirect=%s", cfg.IdentityServerUrl, cfg.UrlRoot)
// 			apiutil.WriteMessageResponse(w, msg, nil)
// 			return
// 		} else if res.StatusCode != http.StatusOK {
// 			log.Info(data["meta"].(map[string]interface{})["message"])
// 			apiutil.WriteErrResponse(w, http.StatusInternalServerError, fmt.Errorf("%s", data["meta"].(map[string]interface{})["message"]))
// 			return
// 		}

// 		perm := data["data"].(map[string]interface{})["permission"]
// 		if perm != "admin" && perm != "write" {
// 			apiutil.WriteErrResponse(w, http.StatusUnauthorized, fmt.Errorf("access denied"))
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "permission", perm)

// 		// no-auth middware func
// 		middleware(handler)(w, r.WithContext(ctx))
// 	}
// }

// addCORSHeaders adds CORS header info for whitelisted servers
func addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	// origin := r.Header.Get("Origin")
	// for _, o := range cfg.AllowedOrigins {
	// 	if origin == o {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// return
	// 	}
	// }
}
