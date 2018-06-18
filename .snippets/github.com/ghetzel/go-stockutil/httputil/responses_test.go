package httputil

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type testOutputOne struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Count int64
	Ok1   bool `json:"OK1" maputil:"ok1"`
	Ok2   bool `json:"OK2" maputil:"ok2"`
	Ok3   bool `json:"OK3"`
}

func TestParseFormValues(t *testing.T) {
	assert := require.New(t)

	t1 := testOutputOne{
		URL: `http://test`,
	}

	assert.NoError(ParseFormValues(url.Values{
		`name`:  []string{`Tester`},
		`Count`: []string{`42`},
		`ok1`:   []string{`true`},
		`ok2`:   []string{`on`},
		`OK3`:   []string{`yes`},
	}, &t1))

	assert.Equal(`Tester`, t1.Name)
	assert.Equal(`http://test`, t1.URL)
	assert.Equal(int64(42), t1.Count)
	assert.True(t1.Ok1)
	assert.True(t1.Ok2)
	assert.True(t1.Ok3)
}
