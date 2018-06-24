package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type AuthMethod interface {
	auth(req *http.Request)
}

type basicAuthCallback func() (username, password string)

func BasicAuth(username, password string) AuthMethod {
	return basicAuthCallback(func() (string, string) { return username, password })
}

func (cb basicAuthCallback) auth(req *http.Request) {
	username, password := cb()
	req.SetBasicAuth(username, password)
}

func (c *ConfluenceService) sendRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json")

	c.authMethod.auth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusPartialContent:
		res, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return res, nil
	case http.StatusNoContent, http.StatusResetContent:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed")
	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("service is not available (%s)", resp.Status)
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal server error: %s", resp.Status)
	}

	return nil, fmt.Errorf("unknown response status %s", resp.Status)
}
