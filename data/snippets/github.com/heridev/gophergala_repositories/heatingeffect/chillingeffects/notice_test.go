package chillingeffects

import (
	"testing"
)

func TestRequestNoticeRandom(t *testing.T) {
	for i := 1; i <= 10; i++ {
		_, err := RequestNotice(i)
		if err != nil {
			t.Error(err)
		}
	}
}
