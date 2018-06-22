package sls

import (
	"encoding/json"
	"fmt"
)

// Error defines sls error
type Error struct {
	Code      string `json:"errorCode"`
	Message   string `json:"errorMessage"`
	RequestID string `json:"requestID"`
}

// NewClientError new client error
func NewClientError(message string) *Error {
	err := new(Error)
	err.Code = "ClientError"
	err.Message = message
	return err
}

func (e Error) String() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

func (e Error) Error() string {
	return e.String()
}

// Client ...
type Client struct {
	Endpoint        string // IP or hostname of SLS endpoint
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
}

func convert(c *Client, projName string) *LogProject {
	return &LogProject{
		Name:            projName,
		Endpoint:        c.Endpoint,
		AccessKeyID:     c.AccessKeyID,
		AccessKeySecret: c.AccessKeySecret,
		SecurityToken:   c.SecurityToken,
	}
}

// CreateProject create a new loghub project.
func (c *Client) CreateProject(name, description string) (*LogProject, error) {
	type Body struct {
		ProjectName string `json:"projectName"`
		Description string `json:"description"`
	}
	body, err := json.Marshal(Body{
		ProjectName: name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%d", len(body)),
		"Content-Type":      "application/json",
		"Accept-Encoding":   "deflate", // TODO: support lz4
	}

	uri := "/"
	proj := convert(c, name)
	_, err = request(proj, "POST", uri, h, body)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// GetProject ...
func (c *Client) GetProject(name string) (*LogProject, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
	}

	uri := "/"
	proj := convert(c, name)
	_, err := request(proj, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// CheckProjectExist check project exist or not
func (c *Client) CheckProjectExist(name string) (bool, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
	}
	uri := "/"
	proj := convert(c, name)
	_, err := request(proj, "GET", uri, h, nil)
	if err != nil {
		if _, ok := err.(*Error); ok {
			slsErr := err.(*Error)
			if slsErr.Code == "ProjectNotExist" {
				return false, nil
			}
			return false, slsErr
		}
		return false, err
	}
	return true, nil
}

// DeleteProject ...
func (c *Client) DeleteProject(name string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
	}

	proj := convert(c, name)
	uri := "/"
	_, err := request(proj, "DELETE", uri, h, nil)
	if err != nil {
		return err
	}

	return nil
}
