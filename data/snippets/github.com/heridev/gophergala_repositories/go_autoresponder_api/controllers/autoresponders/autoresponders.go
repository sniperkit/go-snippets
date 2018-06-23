package autoresponders

import (
  "net/http"
  "github.com/heridev/go_autoresponder_api/utils"
  "github.com/heridev/go_autoresponder_api/models"
  "github.com/heridev/go_autoresponder_api/database"
)

func Index(w http.ResponseWriter, req *http.Request) {
  db := database.DbInstance

  var autoresponders []models.Autoresponder

  db.Find(&autoresponders)

  for index, _ := range autoresponders {
    var lists []models.EmailList
    db.Model(&autoresponders[index]).Related(&lists)
    autoresponders[index].Lists = lists
  }

  if autoresponders == nil {
    utils.Render().JSON(w, http.StatusOK, nil)
  } else {
    utils.Render().JSON(w, http.StatusOK, &autoresponders)
  }
}

func Create(w http.ResponseWriter, req *http.Request) {
}

