package envtokensource

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
)

/*
OAuth2 Utilities
*/

// EnvTokenSource ...
// TODO(alvivi): doc this
type EnvTokenSource struct {
	token *oauth2.Token
}

// NewEnvTokenSource ...
// TODO(alvivi): doc this
func NewEnvTokenSource(envVarName string) (*EnvTokenSource, error) {
	accessToken := os.Getenv(envVarName)
	if len(accessToken) <= 0 {
		return nil, errors.New("No envvar found")
	}
	ts := EnvTokenSource{
		token: &oauth2.Token{
			AccessToken: accessToken,
		},
	}
	return &ts, nil
}

// Token ...
// TODO(alvivi): doc this
func (ts EnvTokenSource) Token() (*oauth2.Token, error) {
	return ts.token, nil
}
