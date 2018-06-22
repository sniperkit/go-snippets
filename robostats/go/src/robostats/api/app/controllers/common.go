package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"io/ioutil"
	"log"
	"net/http"
	"robostats/models/user"
	"strings"
)

type errEnvelope struct {
	Errors []string `json:"errors"`
}

type Common struct {
	*revel.Controller
}

func (c Common) requireAuthorization() (*user.User, error) {
	authHeader := c.Request.Header.Get("Authorization")
	log.Printf("Got auth header: %v\n", authHeader)

	chunks := strings.Split(authHeader, " ")

	token := chunks[len(chunks)-1]
	log.Printf("got token: %x\n", token)

	return user.GetByToken(token)
}

func (c Common) decodeBody(dest interface{}) error {
	var buf []byte
	var err error

	if buf, err = ioutil.ReadAll(c.Request.Body); err != nil {
		return err
	}

	log.Printf("got: %s\n", string(buf))

	if err := json.Unmarshal(buf, &dest); err != nil {
		return err
	}

	return nil
}

func (c Common) writeError(err error) revel.Result {
	c.Response.Status = 422
	c.Response.ContentType = "application/json"
	data := errEnvelope{[]string{err.Error()}}
	log.Printf("app error: %q", err)
	return c.RenderJson(data)
}

func (c Common) writeStatus() revel.Result {
	c.Response.ContentType = "text/plain"
	return c.RenderText(http.StatusText(c.Response.Status) + "\n")
}

func (c Common) dataCreated(data interface{}) revel.Result {
	c.Response.Status = http.StatusCreated
	c.Response.ContentType = "application/json"
	return c.RenderJson(data)
}

func (c Common) dataGeneric(data interface{}) revel.Result {
	c.Response.Status = http.StatusOK
	c.Response.ContentType = "application/json"
	return c.RenderJson(data)
}

func (c Common) StatusOK() revel.Result {
	c.Response.Status = http.StatusOK
	return c.writeStatus()
}

func (c Common) StatusUnauthorized() revel.Result {
	c.Response.Status = http.StatusUnauthorized
	return c.writeStatus()
}

func (c Common) StatusBadRequest() revel.Result {
	c.Response.Status = http.StatusBadRequest
	return c.writeStatus()
}

func (c Common) statusNotFound() revel.Result {
	c.Response.Status = http.StatusNotFound
	return c.writeStatus()
}
