package encryption

type RandomTokenGenerator interface {
	AlphaNumericString(length int) string
	AlphaString(length int) string
	NumericString(length int) string
	CustomString(length int, charlist string) string
}
