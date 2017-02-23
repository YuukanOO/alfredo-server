package domain

import (
	"github.com/YuukanOO/go-toolbelt/errors"
	"github.com/YuukanOO/go-toolbelt/validation"
)

const (
	// ErrUserAlreadyRegistered is thrown when a user has already been registered.
	ErrUserAlreadyRegistered = "UserAlreadyRegistered"
)

// Register a new user in the system.
func Register(userRepo UserRepository, crypto CryptoService, email, password string) (*User, error) {

	if err := validation.Validate("User").
		Field("email", email, "required,email").
		Field("password", password, "required").
		Errors(); err != nil {
		return nil, err
	}

	if usr, _ := userRepo.FindByEmail(email); usr != nil {
		return nil, errors.NewDomainError(ErrUserAlreadyRegistered, "User already registered")
	}

	cryptedPassword, err := crypto.EncryptPassword(password)

	if err != nil {
		return nil, err
	}

	return newUser(email, cryptedPassword), nil
}
