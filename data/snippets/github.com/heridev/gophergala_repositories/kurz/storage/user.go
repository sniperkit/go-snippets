package storage

type User struct {
	Id              int64
	DefaultStrategy string // The name of the user's default strategy
}

func CurrentUser() User {
	var ret User
	return ret
}
