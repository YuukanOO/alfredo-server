package identity

// RegisterUser register a user in alfredo.
func RegisterUser(
	userRepo UserRepository,
	cryptoService CryptoService,
	email string, password string) (*User, error) {

	// Validates unicity of the user email.
	user, err := userRepo.FindByEmail(email)

	if user != nil {
		return nil, err
	}

	encryptedPassword, err := cryptoService.EncryptPassword(password)

	if err != nil {
		return nil, err
	}

	// Creates the user
	user = newUser(email, encryptedPassword)

	return user, nil
}
