package encryption

type Password interface {
	GetHashedPasswordHex() string
	GetSaltHex() string
}
