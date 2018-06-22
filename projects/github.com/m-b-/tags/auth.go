package main

import (
	"errors"
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"net/http"
	"strings"

//	"log"
)

const (
	authserver	=	"https://localhost:8080/"
	key = "yaki6m6XctAFmKjJupcPkWAEbFSvpxNHifch4Nyx44aAnAzBBYtZ8UqKbzRUqJnr"
	cname = "tags-token"

)

var (
	Client = &http.Client{}

	hashKey = []byte(securecookie.GenerateRandomKey(32))
	blockKey = []byte(securecookie.GenerateRandomKey(32))
	s = securecookie.New(hashKey, blockKey)
)

// Clean name for safe mkdir
func cleanName(name string) string {
	return strings.Join(strings.FieldsFunc(name, func (r rune) bool {
		return r == '/' || r == '.'
	}), "")
}

func mkr(descr string) (string, error) {
	resp, err := Client.Get(authserver+"/api/"+descr+"&key="+key)
	// XXX make sure err doesn't embed sensible data (eg. key…)
	if err != nil { return "", err }

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { return "", err }

	return string(body), nil
}

func update() {
}

func login2(login string) (string, error) {
	return mkr("login?login="+login)
}

func chain(token string) (string, error) {
	return mkr("chain?token="+token)
}

func info(token string) (string, error) {
	return mkr("info?token="+token)
}

func logout2(token string) {
	mkr("logout?token="+token)
}

// retrieve token; chain it to the auth server
// return data stored in cookie or err
func ChainToken(w http.ResponseWriter, r *http.Request) (int32, error) {
	token, uid, err := getToken(r)
	if err == nil { token, err = chain(token) }
	if err != nil { LogError(err); return 0, err }

	if token == "ko" {
		return 0, errors.New("bad token")
	}

	// previous token was valid, set new token
	err = setToken(w, token, uid)

	return uid, err
}

func setToken(w http.ResponseWriter, token string, uid int32) error {
	value := map[string]interface{}{
		"token"	:	token,
		"uid"	:	uid,
     	}
  
	if encoded, err := s.Encode(cname, value); err == nil {
		cookie := &http.Cookie{
			Name:	cname,
			Value:	encoded,
			Path:	"/",
		}
		http.SetCookie(w, cookie)
	} else {
		return err
	}

	return nil
}

func unsetToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:	cname,
		Path:	"/",
		MaxAge:	-1,
	}
	http.SetCookie(w, cookie)
}

func getToken(r *http.Request) (token string, uid int32, err error) {
	if cookie, err := r.Cookie(cname); err == nil {
		value := map[string]interface{}{}

		if err = s.Decode(cname, cookie.Value, &value); err == nil {
			return value["token"].(string), value["uid"].(int32), nil
		}
	}

	return "", 0, err
}
