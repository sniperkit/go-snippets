package main

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type User struct {
	UserID      string    `db:"user_id"`
	AccessToken string    `db:"access_token"`
	AvatarURL   string    `db:"avatar_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func initOauth2() error {

	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	context.oauth = conf

	return nil
}

func getRedirectURL() string {
	url := context.oauth.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Println("Redirect URL: ", url)

	return url
}

func createUser(token string) (map[string]string, error) {

	url := strings.Join([]string{"https://api.github.com/user?", "access_token=", token}, "")
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var ghUser map[string]interface{}
	err = json.Unmarshal(body, &ghUser)
	if err != nil {
		log.Panic(err)
	}

	// Insert into DB
	ghUserInt := int(ghUser["id"].(float64))
	ghUserStr := strconv.Itoa(ghUserInt)

	user := User{}
	user.UserID = ghUserStr
	user.AvatarURL = ghUser["avatar_url"].(string)
	user.AccessToken = token
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = context.db.NamedExec("INSERT OR REPLACE INTO user_table (user_id, access_token, avatar_url, created_at, updated_at) VALUES (:user_id, :access_token, :avatar_url, :created_at, :updated_at)", &user)

	//_, err = context.db.Exec(`INSERT INTO user_table (user_id, access_token) VALUES ("aaaa","asdkjasldjas")`)

	if err != nil {
		log.Panic(err)
	}
	var id int
	_ = context.db.Get(&id, "SELECT id FROM user_table WHERE user_id = ?", ghUserStr)

	log.Println("Git User: ", ghUser)

	strid := strconv.Itoa(id)
	us := map[string]string{
		"userid":      strid,
		"github_user": ghUser["login"].(string),
		"avatar_url":  ghUser["avatar_url"].(string),
	}

	return us, nil
}

// Auth - Endpoint
func Auth(w http.ResponseWriter, r *http.Request) {
	initOauth2()
	url := getRedirectURL()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback - Endpoint
// Auth
func Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	log.Println("Code: ", code)

	// Swap temp token
	tok, err := context.oauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Token: ", tok)

	us, _ := createUser(tok.AccessToken)

	// Create session
	session, _ := store.Get(r, "groupify")

	// Set some session values.
	session.Values["userID"] = us["userid"]
	session.Values["github_user"] = us["github_user"]
	session.Values["avatar_url"] = us["avatar_url"]

	// Save it.
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
