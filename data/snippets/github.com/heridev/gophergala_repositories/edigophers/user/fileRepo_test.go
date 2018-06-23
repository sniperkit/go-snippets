package user

import (
	"edigophers/utils"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestUnMarshalUserArray(t *testing.T) {
	var m []User

	body := []byte(`[{"Name":"Guillaume","Location":{"Longitude":-3.34246233,"Latitude":55.86122317},"Interests":[{"Name":"Rugby","Rating":5},{"Name":"Tennis","Rating":8}]},{"Name":"Martin","Location":{"Longitude":-2.8197334,"Latitude":56.017244},"Interests":[{"Name":"Skiing","Rating":7},{"Name":"Cinema","Rating":8.5},{"Name":"Salsa","Rating":3}]}]`)

	err := json.Unmarshal(body, &m)

	if err != nil {
		t.Errorf("(%v)", err)
	}

	log.Printf("(%v)", m)
}

func TestCanReadUserDataFromFile(t *testing.T) {
	repo, err := NewRepo("../data/users.json")
	verifyRepoCreated(t, err)
	users, err := repo.GetUsers()

	if err != nil {
		t.Errorf("GetRecommendations failed: %s", err)
	}

	if len(users) == 0 {
		t.Error("Failed to read users data expected not empty but was empty failed")
	}

	log.Printf("(%v)", users[0])
}

func TestCanSerializeToJson(t *testing.T) {
	var users = []User{
		User{Name: "Guillaume",
			Location: Location{Latitude: 55.86122317, Longitude: -3.34246233},
			Interests: Interests{
				Interest{Name: "Rugby", Rating: 5},
				Interest{Name: "Tennis", Rating: 8}}},
		User{Name: "Martin",
			Location: Location{Latitude: 56.017244, Longitude: -2.8197334},
			Interests: Interests{
				Interest{Name: "Skiing", Rating: 7},
				Interest{Name: "Cinema", Rating: 8.5},
				Interest{Name: "Salsa", Rating: 3}}}}

	_, err := json.Marshal(users)
	if err != nil {
		t.Error(err)
	}
}

func verifyRepoCreated(t *testing.T, err error) {
	if err != nil {
		t.Error("Failed to create repostory!", err)
	}
}

func createTestUser() *User {
	return &User{
		Name:     "Test",
		Location: Location{Latitude: 10, Longitude: -10},
		Interests: Interests{
			Interest{Name: "Board games", Rating: 6}}}
}

func TestCanSaveNewUser(t *testing.T) {

	fileName := "test/users.json"
	f, err := os.Create(fileName)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(fileName)
	defer f.Close()

	usr := createTestUser()
	repo, err := NewRepo(fileName)
	verifyRepoCreated(t, err)

	err = repo.SaveUser(*usr)
	if err != nil {
		t.Error(err)
	}

	repo, err = NewRepo(fileName)
	utils.CheckErrorMsg(err, "Failed to create repository")
	repoUsr, err := repo.GetUser(usr.Name)
	if err != nil {
		t.Error(err)
	}

	if loc := repoUsr.Location; loc.Longitude == 0 {
		t.Errorf("Repo user location is empty: (%v) user: (%v) Expected: (%v)", loc, repoUsr, usr.Location)
	}

	if interests := repoUsr.Interests; len(interests) == 0 || interests[0].Name != usr.Interests[0].Name {
		t.Errorf("Repo user interest is empty: (%v) user: (%v) Expected: (%v)", interests, repoUsr, usr.Interests[0].Name)
	}

}

func TestEditSaveNewUser(t *testing.T) {
	fileName := "test/users2.json"
	f, err := os.Create(fileName)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(fileName)
	defer f.Close()

	usr := createTestUser()

	repo, err := NewRepo(fileName)
	verifyRepoCreated(t, err)
	err = repo.SaveUser(*usr)
	if err != nil {
		t.Error(err)
	}

	repo, err = NewRepo(fileName)
	utils.CheckErrorMsg(err, "Failed to create repository")
	usr, err = repo.GetUser(usr.Name)
	if err != nil {
		t.Error(err)
	}

	usr.Location.Latitude = 15
	usr.Interests = append(usr.Interests, Interest{Name: "Salsa", Rating: 3})

	err = repo.SaveUser(*usr)
	if err != nil {
		t.Error(err)
	}

	repo, err = NewRepo(fileName)
	utils.CheckErrorMsg(err, "Failed to create repository")
	repoUsr, err := repo.GetUser(usr.Name)
	if err != nil {
		t.Error(err)
	}

	if loc := repoUsr.Location; loc.Longitude != usr.Location.Longitude {
		t.Errorf("Repo user location is empty: (%v) user: (%v) Expected: (%v)", loc, repoUsr, usr.Location)
	}

	if interests := repoUsr.Interests; len(interests) != 2 || interests[1].Name != usr.Interests[1].Name {
		t.Errorf("Repo user interest is empty: (%v) user: (%v) Expected: (%v)", interests, repoUsr, usr.Interests[1].Name)
	}
}
