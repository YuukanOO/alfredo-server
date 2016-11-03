package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Room represents a house room.
type Room struct {
	ID        bson.ObjectId `bson:"_id" json:"id" validate:"required"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	CreatedBy bson.ObjectId `bson:"created_by" json:"created_by"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

func newRoom(name string, controller bson.ObjectId) *Room {
	return &Room{
		ID:        bson.NewObjectId(),
		Name:      name,
		CreatedBy: controller,
		CreatedAt: time.Now(),
	}
}

// RegisterDevice registers a device for this room.
func (room *Room) RegisterDevice(
	findDevices QueryFunc,
	name string,
	adapter *Adapter,
	config map[string]interface{}) (*Device, error) {

	var existingRoom Room

	if err := findDevices(ByName(name)).One(&existingRoom); err == nil {
		return nil, ErrDeviceNameAlreadyExists
	}

	// First, validates the config by looking each needed adapter config values
	// in the given config map
	if err := adapter.validateConfig(config); err != nil {
		return nil, err
	}

	device := newDevice(room.ID, name, adapter.ID, config)

	if err := validate.Struct(device); err != nil {
		return nil, err
	}

	return device, nil
}

// Rename a room and check for duplicates.
func (room *Room) Rename(findRooms QueryFunc, newName string) error {
	if err := findRooms(ByName(newName)); err == nil {
		return ErrRoomNameAlreadyExists
	}

	oldName := room.Name
	room.Name = newName

	if err := validate.Struct(room); err != nil {
		room.Name = oldName
		return err
	}

	return nil
}
