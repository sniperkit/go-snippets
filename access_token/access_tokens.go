package access_token

import (
	"github.com/datatogether/sqlutil"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var alphaNumericRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alphaNumericRunes[rand.Intn(len(alphaNumericRunes))]
	}
	return string(b)
}

// Create a new access token, making sure it doesn't already exist
func Create(db sqlutil.Execable) (string, error) {
	var token string
	exists := true
	for exists {
		token = RandString(25)
		if err := db.QueryRow("SELECT exists(SELECT 1 FROM users WHERE access_token=$1)", token).Scan(&exists); err != nil {
			// return "", Error500IfErr(err)
			return "", err
		}
	}

	return token, nil
}
