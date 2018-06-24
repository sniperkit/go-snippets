package main

import (
	"encoding/json"
)

// TODO - abandoned session-over-websockets stuff for now b/c
// to get a session we're currently using cookies, which are far easier
// as ajax calls
type CreateUserAct struct {
	ReqAction
	Username string
	Email    string
	Password string
}

func (CreateUserAct) Type() string        { return "SESSION_SIGNUP_REQUEST" }
func (CreateUserAct) SuccessType() string { return "SESSION_SIGNUP_SUCCESS" }
func (CreateUserAct) FailureType() string { return "SESSION_SIGNUP_FAILURE" }

func (CreateUserAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &CreateUserAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}
func (s *CreateUserAct) Exec() (res *ClientResponse) {
	return &ClientResponse{
		Type:      s.SuccessType(),
		RequestId: s.RequestId,
		Schema:    "SEARCH_RESULT_ARRAY",
		// Data:      results,
	}
}

type SessionLoginAct struct {
	ReqAction
	Username string
	Password string
}

func (SessionLoginAct) Type() string        { return "SESSION_LOGIN_REQUEST" }
func (SessionLoginAct) SuccessType() string { return "SESSION_LOGIN_SUCCESS" }
func (SessionLoginAct) FailureType() string { return "SESSION_LOGIN_FAILURE" }

func (SessionLoginAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SessionLoginAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}
func (s *SessionLoginAct) Exec() (res *ClientResponse) {
	return &ClientResponse{
		Type:      s.SuccessType(),
		RequestId: s.RequestId,
		Schema:    "USER",
		// Data:      results,
	}
}

type SessionLogoutAct struct {
	ReqAction
	Query    string
	Page     int
	PageSize int
}

func (SessionLogoutAct) Type() string        { return "SESSION_LOGOUT_REQUEST" }
func (SessionLogoutAct) SuccessType() string { return "SESSION_LOGOUT_SUCCESS" }
func (SessionLogoutAct) FailureType() string { return "SESSION_LOGOUT_FAILURE" }

func (SessionLogoutAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SessionLogoutAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}
func (s *SessionLogoutAct) Exec() (res *ClientResponse) {
	return &ClientResponse{
		Type:      s.SuccessType(),
		RequestId: s.RequestId,
		Schema:    "USER",
		// Data:      results,
	}
}

type SessionKeysAct struct {
	ReqAction
	Query    string
	Page     int
	PageSize int
}

func (SessionKeysAct) Type() string        { return "SESSION_KEYS_REQUEST" }
func (SessionKeysAct) SuccessType() string { return "SESSION_KEYS_SUCCESS" }
func (SessionKeysAct) FailureType() string { return "SESSION_KEYS_FAILURE" }

func (SessionKeysAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SessionKeysAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}
func (s *SessionKeysAct) Exec() (res *ClientResponse) {
	return &ClientResponse{
		Type:      s.SuccessType(),
		RequestId: s.RequestId,
		// Schema:    "KEY_ARRAY",
		// Data:      results,
	}
}

type SaveUserAct struct {
	ReqAction
	Query    string
	Page     int
	PageSize int
}

func (SaveUserAct) Type() string        { return "SAVE_SESSION_USER_REQUEST" }
func (SaveUserAct) SuccessType() string { return "SAVE_SESSION_USER_SUCCESS" }
func (SaveUserAct) FailureType() string { return "SAVE_SESSION_USER_FAILURE" }

func (SaveUserAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SaveUserAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}
func (s *SaveUserAct) Exec() (res *ClientResponse) {
	return &ClientResponse{
		Type:      s.SuccessType(),
		RequestId: s.RequestId,
		Schema:    "USER",
		// Data:      results,
	}
}
