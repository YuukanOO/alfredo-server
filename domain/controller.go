package domain

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// Controller represents a remote controller such as a phone.
type Controller struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	UID   string        `bson:"uid" json:"-"`
	Token string        `bson:"token" json:"-"`
}

func newController(uid string) *Controller {
	return &Controller{
		ID:  bson.NewObjectId(),
		UID: uid,
	}
}

// CreateRoom tries to create a new room for this controller.
func (c *Controller) CreateRoom(findRooms QueryFunc, name string) (*Room, error) {
	var existingRoom Room

	if err := findRooms(ByName(name)).One(&existingRoom); err == nil {
		return nil, ErrRoomNameAlreadyExists
	}

	room := newRoom(name, c.ID)

	return room, nil
}

// RegisterController registers a controller by its uid. It will returns a valid controller
// with a generated token ready to be used.
func RegisterController(
	findControllers QueryFunc,
	secret []byte,
	uid string) (*Controller, error) {

	var existingController Controller

	// If it already exists, returns the token
	if err := findControllers(ByUID(uid)).One(&existingController); err == nil {
		return &existingController, nil
	}

	controller := newController(uid)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": controller.ID,
	})

	tokenStr, err := token.SignedString(secret)

	if err != nil {
		return nil, err
	}

	controller.Token = tokenStr

	return controller, nil
}
