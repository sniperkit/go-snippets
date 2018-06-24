package main

// render the view to create a reset token
// func CreateResetTokenHandler(w http.ResponseWriter, r *http.Request) {
//  logger.Println(r.FormValue("email"))
//  go func() {
//    t, err := CreateResetToken(appDB, r.FormValue("email"))
//    if err != nil {
//      logger.Println(err)
//      return
//    }

//    // if err := sendPasswordResetEmail(t); err != nil {
//    //  logger.Println(err)
//    //  return
//    // }
//  }()
//  // render(ctx, w, r, "user/resetTokenSent.html", nil, nil)
//  Res(w, data)
// }

// reset a user's password
// func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
//  t := &ResetToken{Id: stringParam(r.Context(), "id")}
//  if err := t.Read(appDB); err != nil {
//    ErrRes(w, err)
//    return
//  }

//  u, err := t.Consume(appDB, r.FormValue("password"))
//  if err != nil {
//    ErrRes(w, err)
//    return
//  }

//  if err := setUserSessionCookie(w, r, u.Id); err != nil {
//    ErrRes(w, err)
//    return
//  }

//  http.Redirect(w, r, "/users/"+u.Username, http.StatusSeeOther)
// }
