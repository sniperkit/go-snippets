package main

import (
	"encoding/json"
	"github.com/datatogether/identity/user"
	"net/http"
)

var UsersRequests = user.UserRequests{Store: appDB}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		ListUsersHandler(w, r)
	case "POST":
		CreateUserHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		SingleUserHandler(w, r)
	case "PUT", "POST":
		SaveUserHandler(w, r)
	// case "DELETE":
	default:
		ErrRes(w, ErrNotFound)
	}
}

func SingleUserHandler(w http.ResponseWriter, r *http.Request) {
	envelope := r.FormValue("envelope") != "false"
	p := &user.UsersGetParams{
		Subject: &user.User{
			Id:       r.FormValue("id"),
			Username: r.FormValue("username"),
			// accessToken: r.FormValue("access_token"),
		},
	}
	res := &user.User{}
	if err := UsersRequests.Get(p, res); err != nil {
		ErrRes(w, err)
		return
	}
	Res(w, envelope, res)
}

// list users or get a single user if supplied with a "username" formValue
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	envelope := r.FormValue("envelope") != "false"
	username := r.FormValue("username")
	id := r.FormValue("id")

	if username != "" || id != "" {
		SingleUserHandler(w, r)
		return
	} else {
		page := PageFromRequest(r)
		p := &user.UsersListParams{
			User:   sessionUser(r),
			Type:   reqUserType(r),
			Limit:  page.Size,
			Offset: page.Offset(),
		}
		res := []*user.User{}
		if err := UsersRequests.List(p, &res); err != nil {
			ErrRes(w, err)
			return
		}
		// TODO - should this be here?
		for _, u := range res {
			// remove sensitive fields
			u.Email = ""
		}

		Res(w, envelope, res)
	}

	// ExecRequest(w, envelope, req)
}

func reqUserType(r *http.Request) user.UserType {
	switch r.FormValue("type") {
	case "community":
		return user.UserTypeCommunity
	case "user":
		return user.UserTypeUser
	default:
		return user.UserTypeNone
	}
}

// Create a user from the api, feed password in as a query param
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// sess := sessionUser(r)
	u := user.NewUser("")
	pw := ""

	if isJsonRequest(r) {
		params := struct {
			Username string
			Email    string
			Password string
		}{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			ErrRes(w, NewFmtError(http.StatusBadRequest, "error decoding json: %s", err.Error()))
			return
		}
		u = &user.User{
			Username: params.Username,
			Email:    params.Email,
			// password: params.Password,
		}
		pw = params.Password
	} else {
		// default to form data values
		u.Username = r.FormValue("username")
		u.Email = r.FormValue("email")
		pw = r.FormValue("password")
	}

	p := &user.UsersCreateParams{
		User:     u,
		Password: pw,
	}

	res := &user.User{}
	if err := UsersRequests.Create(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	// log new user in
	if err := setUserSessionCookie(w, r, u.Id); err != nil {
		ErrRes(w, New500Error(err.Error()))
		return
	}

	Res(w, true, res)
}

// confirm a user's email address
// func UserConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
// 	u := NewUser(stringParam(ctx, "id"))
// 	u.emailConfirmed = true
// 	log.Info(u)
// 	if err := u.Save(appDB); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	if err := AddFlashMessage(w, r, "Thanks! Your email address has been confirmed."); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	http.Redirect(w, r, u.Path(), http.StatusSeeOther)
// }

// SaveUserHandler updates a user
func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	u := &user.User{}
	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			ErrRes(w, NewFmtError(http.StatusBadRequest, err.Error()))
			return
		}
	} else {
		// TODO - fill user out from form values
	}

	p := &user.UsersSaveParams{
		User:    sessionUser(r),
		Subject: u,
	}

	res := &user.User{}
	if err := UsersRequests.Save(p, res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}

func UsersSearchHandler(w http.ResponseWriter, r *http.Request) {
	page := PageFromRequest(r)
	p := &user.UsersSearchParams{
		User:   sessionUser(r),
		Query:  r.FormValue("q"),
		Limit:  page.Size,
		Offset: page.Offset(),
	}

	res := []*user.User{}
	if err := UsersRequests.Search(p, &res); err != nil {
		ErrRes(w, err)
		return
	}
	Res(w, true, res)
}

func UserCommunitiesHandler(w http.ResponseWriter, r *http.Request) {
	page := PageFromRequest(r)
	p := &user.UsersCommunitiesParams{
		User:   &user.User{Id: r.FormValue("id")},
		Order:  "created desc",
		Limit:  page.Size,
		Offset: page.Offset(),
	}

	res := []*user.User{}
	if err := UsersRequests.UserCommunities(p, &res); err != nil {
		ErrRes(w, err)
		return
	}
	Res(w, true, res)
}

func CommunityMembersHandler(w http.ResponseWriter, r *http.Request) {
	page := PageFromRequest(r)
	p := &user.UsersCommunityMembersParams{
		User:      sessionUser(r),
		Community: &user.User{Id: r.FormValue("id")},
		Limit:     page.Size,
		Offset:    page.Offset(),
	}

	res := []*user.User{}
	if err := UsersRequests.CommunityMembers(p, &res); err != nil {
		ErrRes(w, err)
		return
	}

	Res(w, true, res)
}

// delete a user
// func DeleteCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
// 	u := sessionUser(ctx)
// 	if err := u.Delete(appDB); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	session, err := sessionStore.Get(r, "qri.io.user")
// 	if err != nil {
// 		ErrRes(w, err)
// 		return
// 	}
// 	session.Values["id"] = nil
// 	if err := session.Save(r, w); err != nil {
// 		ErrRes(w, err)
// 		return
// 	}

// 	AddFlashMessage(w, r, "You've successfully deleted your account.")
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
