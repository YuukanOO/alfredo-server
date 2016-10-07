package domain

import "gopkg.in/mgo.v2/bson"

// Device represents a smart device connected to our system
type Device struct {
	ID      bson.ObjectId          `bson:"_id" json:"id"`
	RoomID  bson.ObjectId          `bson:"room_id" json:"room_id"`
	Name    string                 `bson:"name" json:"name"`
	Adapter string                 `bson:"adapter" json:"adapter"`
	Config  map[string]interface{} `bson:"config" json:"config"`
}

func newDevice(
	room bson.ObjectId,
	name string,
	adapter string,
	config map[string]interface{}) *Device {
	return &Device{
		ID:      bson.NewObjectId(),
		RoomID:  room,
		Name:    name,
		Adapter: adapter,
		Config:  config,
	}
}
