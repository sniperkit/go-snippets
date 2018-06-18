package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
 )

 func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	 req, _ := http.NewRequest(method, path, nil)
	 w := httptest.NewRecorder()
	 r.ServeHTTP(w, req)

	 return w
 }

 func TestHelloWorld(t *testing.T) {
	 // build our expected body
	 body := gin.H{
		 "hi": "jiangew",
	 }

	 // grab our router
	 router := setupRouter()

	 // perform a GET request with that handler
	 w := performRequest(router, "GET", "/")

	 // assert we encoded correctly, the request gives a 200
	 assert.Equal(t, http.StatusOK, w.Code)

	 // convert the JSON response to a map
	 var response map[string]string
	 err := json.Unmarshal([]byte(w.Body.String()), &response)

	 // grab the value & whether or not it exists
	 value, exists := response["hi"]

	 // make some assertions on the correctness of the response
	 assert.Nil(t, err)
	 assert.True(t, exists)
	 assert.Equal(t, body["hi"], value)
 }
