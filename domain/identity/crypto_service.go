package identity

// CryptoService is a facility to encrypt password in the identity context.
type CryptoService interface {
	// EncryptPassword encrypts the password using a defined algorithm.
	EncryptPassword(password string) (string, error)
	// VerifyPassword checks if a raw password match an hashed one.
	VerifyPassword(hashedPassword string, password string) error
}
