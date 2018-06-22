/*
* @Author: souravray
* @Date:   2015-01-24 10:34:10
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 22:27:47
 */

package controllers

import (
	"fmt"
	"github.com/gophergala/tinyembassy/site/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func LoginPage(rw http.ResponseWriter, req *http.Request) {
	if IsAuth(req) {
		http.Redirect(rw, req, "/campaign/create", 301)
	}
	render(rw, "login.html")
	return
}

func Login(rw http.ResponseWriter, req *http.Request) {
	websession, _ := store.Get(req, "pp-session")
	email := req.FormValue("email")
	pass := req.FormValue("pass")
	encryptedPass := pass //TODO encrypt password
	s, err := mgo.Dial(conf.DbURI)
	c := s.DB(conf.DbName).C("advertiser")

	defer s.Close()

	result := models.Advertiser{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		fmt.Printf("user not found...", err)
		http.Redirect(rw, req, "/", 301)
	} else {
		if result.Pass == encryptedPass {
			fmt.Println("advertiser exist, password matched..")
			websession.Values["id"] = result
			websession.Save(req, rw)
			fmt.Println(rw)
			http.Redirect(rw, req, "/campaign/create", 301)
		} else {
			fmt.Println("password does not match", result.Pass, " == ", encryptedPass)
			http.Redirect(rw, req, "/", 301)
		}
	}
	return
}

func SignupPage(rw http.ResponseWriter, req *http.Request) {
	if IsAuth(req) {
		http.Redirect(rw, req, "/campaign/create", 301)
	}
	render(rw, "signup.html")
	return
}

func Signup(rw http.ResponseWriter, req *http.Request) {
	email := req.FormValue("email")
	pass := req.FormValue("pass")
	name := req.FormValue("name")

	ecrypedPasswd := pass  //TODO: encrypt password
	emailVerified := false //TODO: Implement later..

	s, err := mgo.Dial(conf.DbURI)
	if err != nil {
		panic(err)
	}
	c := s.DB(conf.DbName).C("advertiser")

	defer s.Close()

	result := models.Advertiser{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		cred := models.Advertiser{Id: bson.NewObjectId(), Email: email, Name: name, EmailVerified: emailVerified, Pass: ecrypedPasswd}
		err = cred.Validator()
		if err != nil {
			fmt.Printf("Validation failed: %v\n", err)
			http.Redirect(rw, req, "/", 301)
		}

		err = c.Insert(cred)
		if err != nil {
			fmt.Printf("Can't insert document: %v\n", err)
			http.Redirect(rw, req, "/", 301)
		}

		// Set some session values.
		websession, _ := store.Get(req, "pp-session")
		websession.Values["id"] = cred
		websession.Save(req, rw)
		fmt.Println(websession)
		http.Redirect(rw, req, "/campaign/create", 301)

	} else {
		fmt.Println("advertiser already exist..")
		fmt.Println(result)
		http.Redirect(rw, req, "/", 301)
	}
	return
}

func Logout(rw http.ResponseWriter, req *http.Request) {
	websession, _ := store.Get(req, "pp-session")
	websession.Options.MaxAge = -1
	websession.Save(req, rw)
	http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
}
