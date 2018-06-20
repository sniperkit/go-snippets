package models

type User struct {
	UUID           string   `json:"uuid"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Age            string   `json:"age"`
	Email          string   `json:"email"`
	City           string   `json:"city"`
	State          string   `json:"state"`
	Country        string   `json:"country"`
	LoginLocations []string `login_locations`
	MainPhoto      string   `json:"main_photo"`
}

func createUser(user *User) *User {
	return user
}

func deleteUser(UUID string) error {
	return nil
}

func editUser() {

}
