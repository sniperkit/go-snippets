package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func testHttpServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case `GET`:
			RespondJSON(w, map[string]interface{}{
				`path`: req.URL.Path,
			})
		case `POST`:
			var input interface{}

			if !QBool(req, `thing`) {
				RespondJSON(w, nil, http.StatusForbidden)
				return
			}

			if err := ParseJSON(req.Body, &input); err == nil {
				RespondJSON(w, input, http.StatusCreated)
			} else {
				RespondJSON(w, err, http.StatusBadRequest)
			}

		case `PUT`:
			var input interface{}

			if !QBool(req, `thing`) || !QBool(req, `topthing`) {
				RespondJSON(w, nil, http.StatusForbidden)
				return
			}

			if err := ParseJSON(req.Body, &input); err == nil {
				RespondJSON(w, input, http.StatusAccepted)
			} else {
				RespondJSON(w, err, http.StatusBadRequest)
			}

		case `DELETE`:
			RespondJSON(w, nil)
		}
	}))
}

func TestClient(t *testing.T) {
	assert := require.New(t)
	var out map[string]interface{}
	var outS string

	server := testHttpServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	client.SetParam(`topthing`, true)

	assert.NoError(err)
	assert.NotNil(client)

	// GET
	// --------------------------------------------------------------------------------------------
	response, err := client.Get(`/test/path`, nil, nil)
	assert.NoError(err)
	assert.NotNil(response)
	assert.NoError(ParseJSON(response.Body, &out))
	assert.Equal(map[string]interface{}{
		`path`: `/test/path`,
	}, out)

	// POST
	// --------------------------------------------------------------------------------------------
	response, err = client.Post(`/test/path`, `postable`, map[string]interface{}{
		`thing`: true,
	}, nil)

	assert.NoError(err)
	assert.NotNil(response)
	assert.NoError(client.Decode(response.Body, &outS))
	assert.Equal(`postable`, outS)

	// PUT
	// --------------------------------------------------------------------------------------------
	response, err = client.Put(`/test/path`, `puttable`, map[string]interface{}{
		`thing`: true,
	}, nil)

	assert.NoError(err)
	assert.NotNil(response)
	assert.NoError(client.Decode(response.Body, &outS))
	assert.Equal(`puttable`, outS)

	// PUT
	// --------------------------------------------------------------------------------------------
	response, err = client.Delete(`/test/path`, nil, nil)

	assert.NoError(err)
	assert.NotNil(response)
	assert.Equal(http.StatusNoContent, response.StatusCode)
}
