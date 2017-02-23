package domain

// CryptoService crypt and decrypt password.
type CryptoService interface {
	EncryptPassword(rawPassword string) (string, error)
}
