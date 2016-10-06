package domain

import "errors"

var (
	// ErrAdapterCommandNotFound is thrown when a command could not be processed by an adapter.
	ErrAdapterCommandNotFound = errors.New("AdapterCommandNotFound")
	// ErrRoomNameAlreadyExists is thrown when a room already exists with the same name.
	ErrRoomNameAlreadyExists = errors.New("RoomNameAlreadyExists")
	// ErrDeviceNameAlreadyExists is thrown when a device with the same name already exists.
	ErrDeviceNameAlreadyExists = errors.New("DeviceNameAlreadyExists")
	// ErrDeviceConfigInvalid is thrown when a device config is invalid.
	ErrDeviceConfigInvalid = errors.New("DeviceConfigInvalid")
)
