package identity

// UserRepository is the interface used to persist and retrieve users of the system.
type UserRepository interface {
	FindByEmail(email string) (*User, error)
}
