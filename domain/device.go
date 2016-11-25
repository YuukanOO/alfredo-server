package domain

import (
	"encoding/json"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// Device represents a smart device connected to our system
type Device struct {
	ID      bson.ObjectId          `bson:"_id" json:"id" validate:"required"`
	RoomID  bson.ObjectId          `bson:"room_id" json:"room_id" validate:"required"`
	Name    string                 `bson:"name" json:"name" validate:"required"`
	Adapter string                 `bson:"adapter" json:"adapter" validate:"required"`
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

func validateDeviceName(findDevices QueryFunc, roomID bson.ObjectId, name string) error {
	if count, _ := findDevices(And(ByName(name), ByRoomID(roomID))).Count(); count > 0 {
		return errDeviceNameAlreadyExists
	}

	return nil
}

// Rename a device.
func (device *Device) Rename(findDevices QueryFunc, newName string) error {
	if err := validateDeviceName(findDevices, device.RoomID, newName); err != nil {
		return err
	}

	oldName := device.Name
	device.Name = newName

	if err := validate.Struct(device); err != nil {
		device.Name = oldName
		return err
	}

	return nil
}

// UpdateConfig updates the device configuration.
func (device *Device) UpdateConfig(adapter *Adapter, config map[string]interface{}) error {
	if err := adapter.validateConfig(config); err != nil {
		return err
	}

	device.Config = config

	return nil
}

// UpdateStatus updates the device status based on the givne execution result.
func (device *Device) UpdateStatus(result *ExecutionResult) {
	// TODO: Find why I have to do this
	cleanResult := strings.Replace(strings.Trim(result.Out, "'"), "\\", "", -1)

	// If the execution stdout result is empty, don't update the device status
	if cleanResult == "" {
		return
	}

	// For now, we just check if the output could be parsed into a json interface
	if err := json.Unmarshal([]byte(cleanResult), &device.Status); err != nil {
		// If not, we revert to using just the result string as the new device status
		device.Status = cleanResult
	}
}
