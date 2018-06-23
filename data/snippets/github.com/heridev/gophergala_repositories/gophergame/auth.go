package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	neturl "net/url"
	"os"
)

var oauthCfg = &oauth.Config{
	//TODO: put your project's Client Id here.  To be got from https://code.google.com/apis/console
	ClientId: os.Getenv("ClientId"),

	//TODO: put your project's Client Secret value here https://code.google.com/apis/console
	ClientSecret: os.Getenv("ClientSecret"),

	//For Google's oauth2 authentication, use this defined URL
	AuthURL: "https://github.com/login/oauth/authorize",

	//For Google's oauth2 authentication, use this defined URL
	TokenURL: "https://github.com/login/oauth/access_token",

	//To return your oauth2 code, Google will redirect the browser to this page that you have defined
	//TODO: This exact URL should also be added in your Google API console for this project within "API Access"->"Redirect URIs"
	RedirectURL: "http://" + os.Getenv("Host") + "/authorize",

	//This is the 'scope' of the data that you are asking the user's permission to access. For getting user's info, this is the url that Google has defined.
	// Scope: "https://www.googleapis.com/auth/userinfo.profile",
}

//This is the URL that Google has defined so that an authenticated application may obtain the user's info in json format
const profileInfoURL = "https://api.github.com/user"

// Start the authorization process
func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	//Get the Google URL which shows the Authentication page to the user

	code := r.FormValue("code")

	if code == "" {

		url := oauthCfg.AuthCodeURL("")
		// //redirect user to that page
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	var params neturl.Values = neturl.Values(make(map[string][]string))
	params.Set("client_id", oauthCfg.ClientId)
	params.Set("client_secret", oauthCfg.ClientSecret)
	params.Set("code", code)
	resp, err := http.PostForm(oauthCfg.TokenURL, params)
	if err != nil {
		log.Println(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	tmp, _ := neturl.ParseQuery(string(b))
	token := tmp.Get("access_token")
	if token == "" {
		log.Println("Token empty")
	}

	req, err := http.Get(profileInfoURL + "?access_token=" + token)
	if err != nil {
		log.Println("Request Error:", err)
	}
	defer req.Body.Close()

	body, _ := ioutil.ReadAll(req.Body)

	// w.Write(body)
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
	} else {

		username := m["login"].(string)
		id := fmt.Sprint(m["id"])

		LoginUser(w, r, id, username, token)
		err = AddUser(id, username)
		if err != nil {
			log.Print(err)
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func handleLogout(rw http.ResponseWriter, req *http.Request) {
	LogoutUser(rw, req)
	http.Redirect(rw, req, "/", http.StatusFound)
}
