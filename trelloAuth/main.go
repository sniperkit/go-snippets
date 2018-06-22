package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/zaquestion/kit/db/mysql"
)

const (
	createTrelloUser = "CREATE TABLE IF NOT EXISTS `routines`.`trello_user` (id INT NOT NULL AUTO_INCREMENT UNIQUE, email VARCHAR(50), board VARCHAR(30), credentials VARCHAR(100)) PRIMARY KEY id"
	insertTrelloUser = "INSERT INTO trello_user (email, board, credentials) VALUES (?, ?, ?)"
)

type session struct {
	*oauth.Credentials
	board string
}

var (
	oauthClient = oauth.Client{
		TemporaryCredentialRequestURI: "https://trello.com/1/OAuthGetRequestToken",
		ResourceOwnerAuthorizationURI: "https://trello.com/1/OAuthAuthorizeToken",
		TokenRequestURI:               "https://trello.com/1/OAuthGetAccessToken",
	}
)

func main() {
	clientKey := os.Getenv("TRELLO_DEV_KEY")
	clientSecret := os.Getenv("TRELLO_DEV_SECRET")
	port := os.Getenv("PORT")
	mysqlDB := os.Getenv("MYSQL_DB")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASS")

	err := mysql.DBInit(mysql.Config{mysqlHost, mysqlUser, mysqlPass, mysqlDB})
	if err != nil {
		log.Fatal(err)
	}

	_, err = mysql.Exec(createTrelloUser)
	if err != nil {
		log.Fatal(err)
	}

	oauthClient.Credentials = oauth.Credentials{
		Token:  clientKey,
		Secret: clientSecret,
	}
	http.HandleFunc("/auth", TrelloAuth)
	http.HandleFunc("/oauth", OAuthCallback)
	if err := http.ListenAndServe(":"+port, context.ClearHandler(http.DefaultServeMux)); err != nil {
		log.Fatalf("Error listening, %v", err)
	}

}

type oauthClientInfo struct {
	clientKey    string
	clientSecret string
}

type oauthToken struct {
	oauthToken  string
	oauthSecret string
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func TrelloAuth(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "oauth")
	session.Values["board"] = r.URL.Query().Get("board")
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)

	callback := fmt.Sprintf("http://%s/oauth", r.Host)
	query := url.Values{
		"expiration": {"never"},
		"name":       {"RoutinesScrapper"},
	}
	var err error
	creds, err := oauthClient.RequestTemporaryCredentials(nil, callback, query)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	http.Redirect(w, r, oauthClient.AuthorizationURL(creds, nil), 302)
	session.Values["credentials"] = creds
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "oauth")
	tokenCred, _, err := oauthClient.RequestToken(nil, session.Values["credentials"].(*oauth.Credentials), r.FormValue("oauth_verifier"))
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error getting request token").Error(), 500)
		return
	}

	fmt.Println(tokenCred)

	mysql.Exec(insertTrelloUser, "zaquestion@gmail.com", session.Values["board"].(string), tokenCred.Token+":"+tokenCred.Secret)
}
