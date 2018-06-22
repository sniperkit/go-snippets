package handlers_test

import (
	"net/http"
	"testing"

	"github.com/rms1000watt/golang-project-template/handlers"
	"github.com/stretchr/testify/assert"
)

var PersonHandlerTests = []HandlerTestObj{
	{
		Method: "GET",
		Path:   "/person",
		In: handlers.Person{
			ID: "4571CA38-BF46-4FA7-A507-1D2E124B9F2E",
		},
		Out: handlers.Person{
			ID:   "4571CA38-BF46-4FA7-A507-1D2E124B9F2E",
			Name: "Ryan Smith",
			Age:  201,
		},
		Handler: handlers.PersonHandlerGET,
	},
	{
		Method: "POST",
		Path:   "/person",
		In: handlers.Person{
			ID:   "88888888-4444-4444-4444-BBBBBBBBBBBB",
			Name: "Fish Head",
		},
		Out: handlers.Person{
			ID:   "88888888-4444-4444-4444-BBBBBBBBBBBB",
			Name: "Fish Head",
		},
		Handler: handlers.PersonHandlerPOST,
	},
	{
		Method: "GET",
		Path:   "/person",
		In: handlers.Person{
			ID: "88888888-4444-4444-4444-BBBBBBBBBBBB",
		},
		Out: handlers.Person{
			ID:   "88888888-4444-4444-4444-BBBBBBBBBBBB",
			Name: "Fish Head",
		},
		Handler: handlers.PersonHandlerGET,
	},
}

func TestPersonAll(t *testing.T) {
	for _, personHandlerTest := range PersonHandlerTests {
		w, expectedOut, actualOut := runHandlerTest(t, personHandlerTest)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Bad Code. Error: %s", w.Body.String())
		assert.Equal(t, expectedOut, actualOut, "Not expected output")
	}
}
