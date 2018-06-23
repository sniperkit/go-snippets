package controllers

import (
  "errors"
  "fmt"
  "github.com/gophergala/gotoolbox/models"
  . "github.com/gophergala/gotoolbox/services"
  "github.com/gorilla/mux"
  "github.com/markbates/goth/gothic"
  "net/http"
)

type AuthController struct {
  ApplicationController
}

func GetProviderNameForMux(request *http.Request) (string, error) {
  vars := mux.Vars(request)
  provider := vars["provider"]

  if provider == "" {
    return provider, errors.New("you must select a provider")
  }
  return provider, nil
}

func (controller *AuthController) Create() error {
  gothic.GetProviderName = GetProviderNameForMux

  u, err := gothic.CompleteUserAuth(controller.ResponseWriter, controller.Request)
  if err != nil {
    fmt.Fprintln(controller.ResponseWriter, err)
    return err
  }

  user := models.User{GitHubEmail: u.Email,
    GitHubName:        u.Name,
    GitHubID:          u.UserID,
    GitHubAvatarURL:   u.AvatarURL,
    GitHubAccessToken: u.AccessToken}

  DB().Where(models.User{GitHubID: u.UserID}).FirstOrCreate(&user)

  // @TODO Go idiomatic Id vs. ID
  controller.Session.Values["currentUserId"] = user.Id

  controller.Redirect("/", 200)
  return nil
}

func (controller *AuthController) Delete() error {
  delete(controller.Session.Values, "currentUserId")
  controller.Redirect("/", 200)
  return nil
}
