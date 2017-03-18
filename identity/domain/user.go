package domain

import (
	"github.com/YuukanOO/go-toolbelt/eventsourcing"
	"gopkg.in/mgo.v2/bson"
)

// UserRegistered event when a user is created.
type UserRegistered struct {
	ID       bson.ObjectId
	Email    string
	password string
}

// UserActivated event when a user has been activated.
type UserActivated struct{}

// User is the base resource related to user management.
// For now, a user has full access to alfredo.
type User struct {
	eventsourcing.EventSource

	ID       bson.ObjectId `bson:"_id" json:"id"`
	Email    string        `json:"email"`
	Active   bool          `json:"active"`
	Password string        `json:"-"`
}

func newUser(email, password string) *User {
	usr := &User{}

	eventsourcing.TrackChange(usr, UserRegistered{
		ID:       bson.NewObjectId(),
		Email:    email,
		password: password,
	})

	return usr
}

// Activate this user so he can access and control the house now.
func (u *User) Activate() {
	eventsourcing.TrackChange(u, UserActivated{})
}

// Transition from a state to another.
func (u *User) Transition(evt eventsourcing.Event) {
	switch e := evt.(type) {
	case UserRegistered:
		u.ID = e.ID
		u.Email = e.Email
		u.Password = e.password
		u.Active = false
	case UserActivated:
		u.Active = true
	}
}
