package domain

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// Controller represents a remote controller such as a phone.
type Controller struct {
	ID    bson.ObjectId `bson:"_id"`
	UID   string
	Token string
}

func newController(uid string) *Controller {
	return &Controller{
		ID:  bson.NewObjectId(),
		UID: uid,
	}
}

// CreateRoom tries to create a new room for with this controller.
func (c *Controller) CreateRoom(roomByName RoomByNameFunc, name string) (*Room, error) {
	if _, err := roomByName(name); err == nil {
		return nil, errors.New("RoomNameAlreadyExists")
	}

	room := newRoom(name, c.ID)

	return room, nil
}

// RegisterController registers a controller by its uid. It will returns a valid controller
// with a generated token ready to be used.
func RegisterController(controllerByUID ControllerByUIDFunc, secret []byte, uid string) (*Controller, error) {

	// If it already exists, returns the token
	if existingController, err := controllerByUID(uid); err == nil {
		return existingController, nil
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
