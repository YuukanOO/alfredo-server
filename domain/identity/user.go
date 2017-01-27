package identity

// User as it sounds, represents a user of the system and
// is authorized to access the house.
type User struct {
	Email    string
	password string
}

func newUser(email string, encryptedPassword string) *User {
	return &User{
		Email:    email,
		password: encryptedPassword,
	}
}
