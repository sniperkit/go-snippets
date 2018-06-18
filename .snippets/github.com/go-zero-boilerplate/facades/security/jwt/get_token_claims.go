package jwt

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//GetTokenClaims will get the claims from the JWT token
func GetTokenClaims(token *jwt.Token) (map[string]interface{}, error) {
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims
	jwtClaims, claimsTypeOk := claims.(jwt.MapClaims)
	if !claimsTypeOk {
		return nil, errors.New(fmt.Sprintf("claims type (%T) is wrong", claims))
	}

	return jwtClaims, nil
}
