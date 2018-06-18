package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//GenerateToken will generate a JWT token
func GenerateToken(secretKey string, expiryDuration time.Duration, tokenClaims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)

	for key, val := range tokenClaims {
		claims[key] = val
	}

	if _, ok := claims["exp"]; ok {
		return "", errors.New("cannot set exp from outside jwt.GenerateToken")
	}

	claims["exp"] = time.Now().Add(expiryDuration).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Wrap(err, "unable to sign token")
	}

	return t, nil
}
