package main

type User struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Username    string  `json:"username"`
	AvatarUrl   string  `json:"avatar_url"`
	ProfileUrl  string  `json:"profile_url"`
	Email       string  `json:"email"`
	Joined      string  `json:"joined"`
	Raw         RawJson `json:"raw"`
	AccessToken string  `json:"-"`
	TimeStamp
}

type UserList []User

func (u *UserList) List() error {
	_, err := dbmap.Select(u, "SELECT * FROM USERS ORDER BY CreatedAt DESC")
	return err
}

func (u *User) Create() error {
	// run the UpdateTime ethod on the user model
	u.UpdateTime()

	// run the DB insert function
	err := dbmap.Insert(u)
	return err
}

func (u *User) GetByUsername(n string) error {
	err := dbmap.SelectOne(u, "SELECT * FROM users WHERE username=$1", n)
	return err
}
