package dshash

import (
	"appengine/aetest"
	"bytes"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HandlerWithTestContext struct {
	aetest.Context
}

func (hwtc HandlerWithTestContext) Handle(h WebContextHandler) httprouter.Handle {
	f := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(hwtc.Context, w, r, ps)
	}

	return httprouter.Handle(f)
}

func TestRootRoute(t *testing.T) {
	assert := assert.New(t)

	c, e := aetest.NewContext(nil)
	assert.Nil(e)
	defer c.Close()

	hwtc := HandlerWithTestContext{c}

	router := Router(hwtc)

	r, e := http.NewRequest("GET", "/locations/chischaschos", nil)
	assert.Nil(e)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(w.Code, 200)

	person := Person{}
	e = person.Unmarshal(w.Body.Bytes())

	assert.Nil(e)
	assert.Equal(person.Handler, "chischaschos")
	assert.Empty(person.Locations)

	person.Locations = []string{"Somewhere"}

	personBytes, e := person.Marshal()
	assert.Nil(e)

	r, e = http.NewRequest("POST", "/locations", bytes.NewBuffer(personBytes))
	assert.Nil(e)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(w.Code, 200)

	person.Unmarshal(w.Body.Bytes())

	assert.Nil(e)
	assert.Equal(person.Handler, "chischaschos")
	assert.Equal(person.Locations, []string{"Somewhere"})

}
