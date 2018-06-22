package email_lists

import (
  "net/http"
  "github.com/heridev/go_autoresponder_api/utils"
  "github.com/heridev/go_autoresponder_api/models"
  "github.com/heridev/go_autoresponder_api/database"
)

func IndexHandler(w http.ResponseWriter, req *http.Request) {
  db := database.DbInstance

  var lists []models.EmailList
  db.Find(&lists)
  if lists == nil {
    // If we have no subscribers, just return an empty array, instead of null.
    utils.Render().JSON(w, http.StatusOK, "[]")
  } else {
    utils.Render().JSON(w, http.StatusOK, &lists)
  }
}

