package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/datatogether/identity/jwt"
	"github.com/datatogether/identity/user"
	"net"
	"net/http"
	"strings"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		GetSessionHandler(w, r)
	case "PUT":
		SaveUserHandler(w, r)
	case "POST":
		LoginHandler(w, r)
	case "DELETE":
		LogoutHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

// grab the current serialized sesion user, returns nil if no user present
func sessionUser(r *http.Request) *user.User {
	if u, ok := r.Context().Value("user").(*user.User); ok {
		return u
	}
	return nil
}

func cookieUser(r *http.Request) *user.User {
	if session, err := sessionStore.Get(r, cfg.UserCookieKey); err == nil {
		if session.Values["id"] != nil {
			if id, ok := session.Values["id"].(string); ok {
				return user.NewUser(id)
			}
		}
	} else {
		log.Infoln(err.Error())
	}
	return nil
}

// check if a user has been provided via access_token param either as a header
// or as request params
func tokenUser(r *http.Request) *user.User {
	u := user.NewUser("")
	fmt.Println(r.Header.Get("access_token"))
	if r.Header.Get("access_token") != "" {
		u = user.NewAccessTokenUser(r.Header.Get("access_token"))
		return u
	} else if r.FormValue("access_token") != "" {
		u = user.NewAccessTokenUser(r.FormValue("access_token"))
		return u
	} else {
		return nil
	}
}

// attempt to extract & read a session user from a given request.
// if no user is provided, an anonymous user is created
func userFromRequest(db *sql.DB, r *http.Request) (*user.User, error) {
	var (
		u   *user.User
		err error
	)

	u = tokenUser(r)
	if u == nil {
		u, err = jwt.JwtUser(db, r)
		if err != nil {
			// logger.Println(err.Error())
		}
	}
	if u == nil {
		u = cookieUser(r)
	}
	if u == nil {
		// create anononymous user with ip address
		return &user.User{
			Username:  getIP(r),
			Anonymous: true,
		}, nil
	}

	if err := u.Read(db); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return u, nil
}

func getIP(r *http.Request) string {
	remoteAddr := r.Header.Get("x-forwarded-for")
	if remoteAddr != "" {
		return strings.TrimSpace(strings.Split(remoteAddr, ",")[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}

// set a user's session cookie so we can track them
func setUserSessionCookie(w http.ResponseWriter, r *http.Request, id string) error {
	log.Infoln("set user cookie:", id)
	session, err := sessionStore.Get(r, cfg.UserCookieKey)
	if err != nil {
		log.Infoln("setUserSessionCookie error", err.Error())
		if session != nil {
			// if we still get a session object back
			// clear the cookie, b/c this one clearly doesn't work
			session.Options.MaxAge = -1
		}
	}
	session.Values["id"] = id
	return session.Save(r, w)
}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)

	envelope := r.FormValue("envelope") != "false"
	if u == nil || u.Id == "" {
		ErrRes(w, NewFmtError(http.StatusUnauthorized, "unauthorized"))
		return
	} else {
		Res(w, envelope, u)
	}
}

func SessionCommunitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		GetSessionCommunities(w, r)
	// case "POST":
	// 	CreateUserKeyHandler(w, r)
	// case "DELETE":
	// DeleteUserKeyHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func GetSessionCommunities(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)
	if u == nil || u.Id == "" {
		ErrRes(w, NewFmtError(http.StatusUnauthorized, "unauthorized"))
		return
	}
	page := PageFromRequest(r)

	p := &user.UsersCommunitiesParams{
		User:   u,
		Order:  "created DESC",
		Limit:  page.Size,
		Offset: page.Offset(),
	}
	res := []*user.User{}
	if err := UsersRequests.UserCommunities(p, &res); err != nil {
		ErrRes(w, err)
		return
	}
	Res(w, true, res)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login struct {
		Username string
		Password string
	}

	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			e := NewFmtError(http.StatusBadRequest, "error parsing json: '%s'", err.Error())
			log.Infoln(e)
			ErrRes(w, e)
			return
		}
	} else {
		// default to using form values
		login.Username = r.FormValue("username")
		login.Password = r.FormValue("password")
	}

	u, err := user.AuthenticateUser(appDB, login.Username, login.Password)
	if err != nil {
		log.Infoln(ErrInvalidUserNamePasswordCombo)
		ErrRes(w, ErrInvalidUserNamePasswordCombo)
		return
	}

	if err := setUserSessionCookie(w, r, u.Id); err != nil {
		log.Infoln(err)
		// ErrRes(w, New500Error(err.Error()))
		// return
	}

	log.Infof("user api login: %s", login.Username)
	Res(w, true, u)
}

// logout a user, overwriting their session cookie with ""
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, cfg.UserCookieKey)
	if err != nil {
		log.Infof("session get user error", err.Error())
	} else {
		if id, ok := session.Values["id"].(string); ok {
			u := user.NewUser(id)
			if err := u.Read(appDB); err == nil {
				log.Info("logout user: %s", u.Username)
			}
		}
	}

	// regardless of what happens in relation to errors
	// if we have a session object, clear it.
	if session != nil {
		session.Values["id"] = nil
		session.Options.MaxAge = -1
	}
	if err := session.Save(r, w); err != nil {
		log.Infoln(err.Error())
		ErrRes(w, err)
		return
	}

	MessageResponse(w, "successfully logged out", nil)
}

func AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		GetAccessTokenHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func GetAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)
	if u == nil || u.Id == "" {
		ErrRes(w, NewFmtError(http.StatusUnauthorized, "unauthorized"))
		return
	} else {
		if err := u.ReadApiToken(appDB); err != nil {
			ErrRes(w, NewFmtError(http.StatusInternalServerError, "error reading user: %s", err.Error()))
			return
		}
		Res(w, true, map[string]string{
			"access_token": u.AccessToken(),
		})
	}
}
