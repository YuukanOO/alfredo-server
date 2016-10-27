package domain

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

// Device represents a smart device connected to our system
type Device struct {
	ID      bson.ObjectId          `bson:"_id" json:"id"`
	RoomID  bson.ObjectId          `bson:"room_id" json:"room_id"`
	Name    string                 `bson:"name" json:"name"`
	Adapter string                 `bson:"adapter" json:"adapter"`
	Config  map[string]interface{} `bson:"config" json:"config"`
	Status  interface{}            `bson:"status" json:"status"`
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

// Rename a device.
func (device *Device) Rename(findDevices QueryFunc, newName string) error {
	if err := findDevices(ByName(newName)); err == nil {
		return ErrDeviceNameAlreadyExists
	}

	device.Name = newName

	return nil
}

// UpdateConfig updates the device configuration.
func (device *Device) UpdateConfig(adapter *Adapter, config map[string]interface{}) error {
	if err := adapter.ValidateConfig(config); err != nil {
		return err
	}

	device.Config = config

	return nil
}

// UpdateStatus updates the device status based on the givne execution result.
func (device *Device) UpdateStatus(result *ExecutionResult) {
	// If the execution stdout result is empty, don't update the device status
	if result.Out == "" {
		return
	}

	// For now, we just check if the output could be parsed into a json interface
	if err := json.Unmarshal([]byte(result.Out), &device.Status); err != nil {
		// If not, we revert to using just the result string as the new device status
		device.Status = result.Out
	}
}
