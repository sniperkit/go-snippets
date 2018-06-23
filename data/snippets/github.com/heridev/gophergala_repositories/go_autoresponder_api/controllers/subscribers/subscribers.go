package subscribers

import (
  "net/http"
  "github.com/heridev/go_autoresponder_api/utils"
  "github.com/heridev/go_autoresponder_api/models"
  "github.com/heridev/go_autoresponder_api/database"
)

func IndexHandler(w http.ResponseWriter, req *http.Request) {
  db := database.DbInstance

  var subscribers []models.Subscriber
  db.Find(&subscribers)
  if subscribers == nil {
    // If we have no subscribers, just return an empty array, instead of null.
    utils.Render().JSON(w, http.StatusOK, nil)
  } else {
    utils.Render().JSON(w, http.StatusOK, &subscribers)
  }
}

func Create(w http.ResponseWriter, req *http.Request) {

}
