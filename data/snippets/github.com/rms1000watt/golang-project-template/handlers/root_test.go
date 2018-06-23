package handlers_test

import (
	"net/http"
	"testing"

	"github.com/rms1000watt/golang-project-template/handlers"
	"github.com/stretchr/testify/assert"
)

var RootHandlerTests = []HandlerTestObj{
	{
		Method:  "GET",
		Path:    "/fish",
		Handler: handlers.RootHandler,
	},
	{
		Method:  "POST",
		Path:    "/fish",
		Handler: handlers.RootHandler,
	},
	{
		Method:  "PUT",
		Path:    "/fish",
		Handler: handlers.RootHandler,
	},
	{
		Method:  "OPTIONS",
		Path:    "/fish",
		Handler: handlers.RootHandler,
	},
}

func TestRootAll(t *testing.T) {
	for _, rootHandlerTest := range RootHandlerTests {
		w, _, _ := runHandlerTest(t, rootHandlerTest)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Bad Code. Error: %s", w.Body.String())
	}
}
