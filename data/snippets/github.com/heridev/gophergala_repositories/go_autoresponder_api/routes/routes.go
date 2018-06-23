package routes

import(
  "github.com/heridev/go_autoresponder_api/controllers/subscribers"
  "github.com/heridev/go_autoresponder_api/controllers/email_lists"
  "github.com/heridev/go_autoresponder_api/controllers/autoresponders"
  "github.com/gorilla/mux"
)

func Create() *mux.Router {
  r := mux.NewRouter()
  r.HandleFunc("/subscribers",  subscribers.IndexHandler ).Methods("GET")
  r.HandleFunc("/subscribers",  subscribers.Create ).Methods("POST")
  r.HandleFunc("/lists",  email_lists.IndexHandler).Methods("GET")
  r.HandleFunc("/autoresponders",  autoresponders.Index).Methods("GET")
  r.HandleFunc("/autoresponders",  autoresponders.Create).Methods("POST")
  return r
}
