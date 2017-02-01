package identity

import "gopkg.in/mgo.v2/bson"
import "github.com/YuukanOO/go-toolbelt/eventsourcing"

type UserCreated struct {
	ID       bson.ObjectId
	Email    string
	password string
}

// User as it sounds, represents a user of the system and
// is authorized to access the house.
type User struct {
	eventsourcing.EventSource
	ID       bson.ObjectId
	Email    string
	password string
}

func newUser(email string, encryptedPassword string) *User {
	user := &User{}
	eventsourcing.TrackChange(user, UserCreated{
		ID:       bson.NewObjectId(),
		Email:    email,
		password: encryptedPassword,
	})
	return user
}

func (u *User) Transition(event eventsourcing.Event) {
	switch evt := event.(type) {
	case UserCreated:
		u.ID = evt.ID
		u.Email = evt.Email
		u.password = evt.password
	}
}
