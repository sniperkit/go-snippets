package core

import (
	"encoding/json"
	"log"
	"net/http"
	"user"
)

// This will be our datastore for now.
var users map[string]user.User = make(map[string]user.User, 0)

func submitHandler(w http.ResponseWriter, r *http.Request) {

	var (
		err error
	)

	if (r.Method == "POST") || (r.Method == "PUT") {
		// Unmarhsall the incomming data as a user.
		var user_in user.User
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(user_in)
		if err != nil {
			log.Println("I failed at unmarshelling the incoming JSON.")
		}

		// Is the submission we received valid?
		if user_in.UTCTimeStamp == 0 {
			log.Println("Incomming submission does not include a timestamp.")
			return
		} else if user_in.Username == "" {
			log.Println("Incomming submission does not include a username.")
			return
		}

		// Get the existing submission if available.
		var user_db user.User
		user_db = users[user_in.Username]
		if user_db.UTCTimeStamp > user_in.UTCTimeStamp {
			log.Println("Incomming submission is older than the existing submission.")
			return
		}

		// If we've gotten this far, the submitted user is valid and supercedes our exisiting user.
		users[user_in.Username] = user_in
	}
}
