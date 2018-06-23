package encryption

type Service interface {
	PasswordMatchHex(guessedPassword, hexSalt, hexHashedPassword string) (bool, error)
	CreatePassword(raw_pass string) (Password, error)
}
