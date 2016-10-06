package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Room represents a house room.
type Room struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `json:"name"`
	CreatedBy bson.ObjectId `json:"created_by"`
	CreatedAt time.Time     `json:"created_at"`
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
	deviceByName DeviceByNameFunc,
	name string,
	adapter *Adapter,
	config map[string]interface{}) (*Device, error) {

	if _, err := deviceByName(name); err == nil {
		return nil, ErrDeviceNameAlreadyExists
	}

	// First, validates the config by looking each needed adapter config values
	// in the given config map
	for ck := range adapter.Config {
		// TODO: type checking maybe...
		if config[ck] == nil {
			return nil, ErrDeviceConfigInvalid
		}
	}

	device := newDevice(room.ID, name, adapter.ID, config)

	return device, nil
}
