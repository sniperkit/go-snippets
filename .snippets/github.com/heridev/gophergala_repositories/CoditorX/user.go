// Copyright (c) 2015, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/sha1"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"crypto/md5"
)

const (
	user_exists         string = "user_exists"
	user_create_failure string = "user_create_failure"
	user_login_error    string = "login_error"

	user_session string = "user_session"
)

type User struct {
	Username string
	Password string
	Email    string
	Locale   string
	Created  int64 // user create time in unix nano
	Updated  int64 // preference update time in unix nano
}

func (u *User) getBasePath() string {
	return filepath.Join(conf.Workspace, u.Username)
}

func (u *User) getFileBasePath(fileName string) string {
	return filepath.Join(conf.Workspace, u.Username, fileName)
}

func getBasePath(username string) string {
	return filepath.Join(conf.Workspace, username)
}

func (u *User) getWorkspace() string {
	return filepath.Join(conf.Workspace, u.Username, "workspace")
}

func getWorkspace(username string) string {
	return filepath.Join(conf.Workspace, username, "workspace")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	model := map[string]interface{}{}
	if "GET" == r.Method {
		toHtml(w, "login.html", model, conf.Locale)
		return
	}

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		model["succ"] = false
		toJson(w, model)
		return
	}

	user.Locale = conf.Locale
	loginerror := login(user, w, r)

	if loginerror != "" {
		model["succ"] = false
		model["msg"] = loginerror
	} else {
		model["succ"] = true
	}
	toJson(w, model)

}

func login(user *User, w http.ResponseWriter, r *http.Request) string {

	// check username
	if _, err := os.Stat(filepath.Join(conf.Workspace, user.Username)); err == nil {

		bytes, err := ioutil.ReadFile(filepath.Join(conf.Workspace, user.Username, "user.json"))
		if nil != err {
			logger.Error(err)
			return getMsg(user.Locale, user_login_error).(string)
		}
		userInfo := &User{}

		err = json.Unmarshal(bytes, userInfo)
		if err != nil {
			logger.Error(err)
			return getMsg(user.Locale, user_login_error).(string)
		}

		if tosha1(user.Password) == userInfo.Password {
			err = saveSession(userInfo, w, r)
			if nil != err {
				logger.Error(err)
				return getMsg(user.Locale, user_login_error).(string)
			}

		} else {
			return getMsg(user.Locale, user_login_error).(string)
		}
	} else {
		return getMsg(user.Locale, user_login_error).(string)
	}
	return ""
}

func saveSession(user *User, w http.ResponseWriter, r *http.Request) error {

	httpSession, _ := httpSessionStore.Get(r, "coditor-session")

	gob.Register(user)
	httpSession.Values[user_session] = user
	httpSession.Values["username"] = user.Username
	httpSession.Values["id"] = strconv.Itoa(rand.Int())

	logger.Trace(httpSession.Values)

	httpSession.Options.MaxAge = conf.HTTPSessionMaxAge
	if "" != conf.Context {
		httpSession.Options.Path = conf.Context
	}

	return httpSession.Save(r, w)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	model := map[string]interface{}{}
	if "GET" == r.Method {
		toHtml(w, "sign_up.html", model, conf.Locale)
		return
	}

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		model["succ"] = false
		model["msg"] = err
		toJson(w, model)
		return
	}

	//TODO
	user.Locale = conf.Locale
	singuperror := signUp(user, w, r)

	if singuperror != "" {

		model["succ"] = false
		model["msg"] = singuperror
	} else {
		model["succ"] = true
	}

	toJson(w, model)
}

func signUp(user *User, w http.ResponseWriter, r *http.Request) string {

	// check if username exists?
	if _, err := os.Stat(filepath.Join(conf.Workspace, user.Username)); err == nil {
		return getMsg(user.Locale, user_exists).(string)
	}

	//create user dir
	userBaseDir := filepath.Join(conf.Workspace, user.Username)
	err := os.MkdirAll(userBaseDir, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return getMsg(conf.Locale, user_create_failure).(string)
	}

	user.Password = tosha1(user.Password)
	now := time.Now().UnixNano()
	user.Created = now
	user.Updated = now

	userjson, _ := json.MarshalIndent(user, "", "    ")

	// create user info json
	fout, err := os.Create(filepath.Join(userBaseDir, "user.json"))
	defer fout.Close()
	if err != nil {
		logger.Error(err)
		return getMsg(user.Locale, user_create_failure).(string)

	}
	_, err = fout.Write(userjson)
	if err != nil {
		logger.Error(err)
		return getMsg(user.Locale, user_create_failure).(string)
	}

	//make user workspace
	err = os.MkdirAll(filepath.Join(userBaseDir, "workspace"), os.ModePerm)
	if err != nil {
		logger.Error(err)
		return getMsg(user.Locale, user_create_failure).(string)

	}

	err = saveSession(user, w, r)
	if nil != err {
		logger.Error(err)
		return getMsg(user.Locale, user_create_failure).(string)
	}
	return ""
}

func tosha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)

	return fmt.Sprintf("%x", t.Sum(nil))
}


func toMd5(data string) string {
	t := md5.New()
	io.WriteString(t, data)

	return fmt.Sprintf("%x", t.Sum(nil))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	model := map[string]interface{}{"succ": true}
	httpSession, _ := httpSessionStore.Get(r, "coditor-session")

	httpSession.Options.MaxAge = -1
	err := httpSession.Save(r, w)
	if err != nil {
		logger.Error(err)
		model["succ"] = false

	}

	toJson(w, model)
}
