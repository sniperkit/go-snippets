package oauth

import (
	"github.com/datatogether/identity/user"
	"golang.org/x/oauth2"
)

const (
	OauthServiceGithub = "github"
)

type OauthUserService interface {
	CurrentUserInfo() (map[string]interface{}, error)
	ExtractUser() (*user.User, error)
}

var (
	GithubOAuth *oauth2.Config
)

func InitOauth(appId, appSecret string) {
	GithubOAuth = &oauth2.Config{
		ClientID:     appId,
		ClientSecret: appSecret,
		// RedirectURL:  "https://ident.archivers.space/oauth/callback",
		Scopes: []string{"user:email", "public_repo"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}
