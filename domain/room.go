package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Room represents a house room.
type Room struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	CreatedBy bson.ObjectId
	CreatedAt time.Time
}

func newRoom(name string, controller bson.ObjectId) *Room {
	return &Room{
		ID:        bson.NewObjectId(),
		Name:      name,
		CreatedBy: controller,
		CreatedAt: time.Now(),
	}
}
