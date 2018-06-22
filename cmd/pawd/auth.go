package main

import (
	"context"
	"errors"

	"github.com/Xe/ln"
	pawd "github.com/Xe/pawd/proto"
)

// Auth is the implementation of the grpc service pawd.Auth.
type Auth struct{ *Server }

func (a *Auth) Register(ctx context.Context, ri *pawd.RegisterInfo) (*pawd.UserToken, error) {
	u, err := a.us.Create(ctx, ri.Email, ri.Password, false)
	if err != nil {
		return nil, err
	}

	ln.Log(ctx, u, ln.Action("new user created"))

	tkn, err := a.tk.Create(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	ln.Log(ctx, tkn, ln.Action("new token created"))

	put := &pawd.UserToken{
		Token: tkn.Body,
	}

	if u.Admin {
		put.Flags = append(put.Flags, "admin")
	}

	return put, nil
}

func (a *Auth) Login(ctx context.Context, li *pawd.LoginInfo) (*pawd.UserToken, error) {
	u, err := a.us.CheckPassword(ctx, li.Email, li.Password, li.TotpChallenge)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid email, password or totp challenge")
	}

	tkn, err := a.tk.Create(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	ln.Log(ctx, u, tkn, ln.Action("new login and token created"))

	return nil, errors.New("not implemented")
}

func (a *Auth) Logout(ctx context.Context, ut *pawd.UserToken) (*pawd.Nil, error) {
	_, err := a.tk.Check(ctx, ut.Token)
	if err != nil {
		return nil, err
	}

	return &pawd.Nil{}, nil
}
