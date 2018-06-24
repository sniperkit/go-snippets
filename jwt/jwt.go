package jwt

import (
	"crypto/rsa"
	"database/sql"
	"github.com/datatogether/identity/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
	"time"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// ArchiversJWTClaims object
type ArchiversClaims struct {
	*jwt.StandardClaims
	UserId   string `json:"userId"`
	Username string `json:"username"`
	UserType string `json:"userType"`
}

// TODO - refactor to a "jwt" object that holds this state
func InitKeys(publicKey, privateKey string) (err error) {
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	return
}

func CreateToken(user *user.User) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims = &ArchiversClaims{
		&jwt.StandardClaims{
			// set the expire time
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
			ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
		},
		user.Id,
		user.Username,
		user.Type.String(),
	}

	// Creat token string
	return t.SignedString(signKey)
}

func JwtUser(db *sql.DB, r *http.Request) (*user.User, error) {
	// Get token from request
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &ArchiversClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})

	// If the token is missing or invalid, return error
	if err != nil {
		return nil, err
	}

	// Token is valid
	// fmt.Fprintln(w, "Welcome,", token.Claims.(*ArchiversClaims).Name)
	u := &user.User{Id: token.Claims.(*ArchiversClaims).UserId}
	err = u.Read(db)
	return u, err
}
