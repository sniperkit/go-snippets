package user

import (
	"edigophers/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
)

//FileRepo is a file based user repository
type FileRepo struct {
	filepath string
	users    []User
}

// GetUser gets a user per his username
func (r FileRepo) GetUser(name string) (*User, error) {
	var result *User
	for _, u := range r.users {
		if u.Name == name {
			result = &u
			break
		}
	}

	if result == nil {
		return nil, errors.New("User doesn't exists")
	}

	return result, nil
}

//GetUsers is a function returning a list of users
func (r FileRepo) GetUsers() ([]User, error) {
	return r.users, nil
}

//SaveUser is a function to persist user changes
func (r FileRepo) SaveUser(usr User) error {
	_, err := r.GetUser(usr.Name)

	if err != nil {
		r.users = append(r.users, usr)
	} else {
		for i, val := range r.users {
			if val.Name == usr.Name {
				r.users[i] = usr
			}
		}
	}

	err = r.saveToFile()
	if err != nil {
		return err
	}

	return nil
}

func (r FileRepo) saveToFile() error {

	body, err := json.Marshal(r.users)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(r.filepath, body, 0600)
	if err != nil {
		return err
	}

	return nil
}

func loadFile(filepath string) ([]User, error) {

	body, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	var usersJSON []User
	if len(body) > 0 {
		err = json.Unmarshal(body, &usersJSON)
		if err != nil {
			return nil, err
		}
	}

	return usersJSON, nil
}

//NewRepo creates a new File base repository
func NewRepo(filepath string) (Repository, error) {

	u, err := loadFile(filepath)
	if err != nil {
		return nil, err
	}
	utils.CheckError(err)
	repo := FileRepo{users: u, filepath: filepath}
	return Repository(repo), nil
}
