// Package cookies provides a convenience wrapper around gorilla/securecookie with shared hash keys
package cookies

import (
	"net/http"
	"os"

	"github.com/bradberger/context"
	"github.com/bradberger/rest"

	"github.com/gorilla/securecookie"
)

var (
	jar *securecookie.SecureCookie
)

const (
	AuthTokenKey = "token"
)

func init() {
	hashKey, blockKey := os.Getenv("COOKIE_HASH_KEY"), os.Getenv("COOKIE_BLOCK_KEY")
	if hashKey != "" && blockKey != "" {
		Init([]byte(hashKey), []byte(blockKey))
	}
}

func getEnvWithDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// Init initializes the cookie jar and is required if you don't want panics. If the
// COOKIE_HASH_KEY and COOKIE_BLOCK_KEY environment variables are supplied it will
// have already been initialized using those values so this will be unnecessary.
func Init(hashKey, blockKey []byte) {
	jar = securecookie.New(hashKey, blockKey)
}

// Set stores the given value with the given key
func Set(ctx context.Context, key string, value interface{}) error {
	w := rest.ResponseWriter(ctx)
	encoded, err := jar.Encode(key, value)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:  key,
		Value: encoded,
		Path:  "/",
		// HttpOnly: true,
		// Secure:   true,
	})
	return nil
}

// Get retrieves a cookie with the given key and encodes it into value
func Get(ctx context.Context, key string, value interface{}) error {
	r := rest.Request(ctx)
	cookie, err := r.Cookie(key)
	if err != nil {
		return err
	}
	if err := jar.Decode(key, cookie.Value, value); err != nil {
		return err
	}
	return nil
}

func Delete(ctx context.Context, key string) error {
	w := rest.ResponseWriter(ctx)
	http.SetCookie(w, &http.Cookie{
		Name:   key,
		Path:   "/",
		MaxAge: -1,
	})
	return nil
}

func Clear(ctx context.Context) error {
	r := rest.Request(ctx)
	for _, cookie := range r.Cookies() {
		if err := Delete(ctx, cookie.Name); err != nil {
			return err
		}
	}
	return nil
}
