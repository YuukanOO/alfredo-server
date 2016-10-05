package domain

import (
	"bytes"
	"errors"

	"gopkg.in/mgo.v2/bson"
)

// Device represents a smart device connected to our system
type Device struct {
	ID      bson.ObjectId          `bson:"_id" json:"id"`
	RoomID  bson.ObjectId          `json:"room_id"`
	Name    string                 `json:"name"`
	Adapter string                 `json:"adapter"`
	Config  map[string]interface{} `json:"config"`
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

// Execute a command for this device with the given parameters.
func (dev *Device) Execute(adapter *Adapter, command string, params map[string]interface{}) error {

	// TODO: Execute should be in the Adapter class!

	tmpl := adapter.commandsParsed[command]

	if tmpl == nil {
		return errors.New("CommandNotFound")
	}

	ctx := newExecutionContext(dev.Config, params)

	var buf bytes.Buffer

	err := tmpl.Execute(&buf, ctx)

	if err != nil {
		return err
	}

	// TODO: launch the command with buf.String()

	return nil
}
