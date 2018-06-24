package access_token

import (
	"testing"
)

func TestNewAccessToken(t *testing.T) {
	token, err := Create(testDB)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(token) != 25 {
		t.Errorf("token is the wrong length. expected: %d, got: %d", 25, len(token))
	}
}
