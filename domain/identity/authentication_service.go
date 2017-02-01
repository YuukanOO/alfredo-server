package identity

import "github.com/YuukanOO/go-toolbelt/errors"
import "github.com/YuukanOO/go-toolbelt/validation"

const (
	// ErrUserEmailAlreadyExists is thrown when an user email is already registered.
	ErrUserEmailAlreadyExists = "UserEmailAlreadyExists"
	// ErrWrongUserPassword is thrown when an invalid password has been provided.
	ErrWrongUserPassword = "WrongUserPassword"
	// ErrUnknownUser is thrown when a user could not be found.
	ErrUnknownUser = "UnknownUser"
)

// Token is a simple alias to represents an access token.
type Token string

// RegisterUser register a user in alfredo.
func RegisterUser(
	userRepo UserRepository,
	cryptoService CryptoService,
	email string, password string) (*User, error) {

	if err := validation.Validate("User").
		Field("email", email, "required,email").
		Field("password", password, "required").Errors(); err != nil {
		return nil, err
	}

	user, err := userRepo.FindByEmail(email)

	if user != nil {
		return nil, errors.NewDomainError(ErrUserEmailAlreadyExists, "User email already registered")
	}

	encryptedPassword, err := cryptoService.EncryptPassword(password)

	if err != nil {
		return nil, err
	}

	user = newUser(email, encryptedPassword)

	return user, nil
}

// AuthenticateUser tries to authenticate a user.
func AuthenticateUser(
	userRepo UserRepository,
	cryptoService CryptoService,
	tokenService TokenService,
	email string, password string) (Token, error) {

	if err := validation.Validate("User").
		Field("email", email, "required,email").
		Field("password", password, "required").Errors(); err != nil {
		return "", err
	}

	user, err := userRepo.FindByEmail(email)

	if err != nil {
		return "", errors.NewDomainError(ErrUnknownUser, "User not found")
	}

	if err = cryptoService.VerifyPassword(user.password, password); err != nil {
		return "", errors.NewDomainError(ErrWrongUserPassword, "Password is not valid for this user")
	}

	token, err := tokenService.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
