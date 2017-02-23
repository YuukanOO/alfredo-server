package domain

// UserRepository deals with User persistence.
type UserRepository interface {
	Add(user *User) error
	FindByEmail(email string) (*User, error)
}
