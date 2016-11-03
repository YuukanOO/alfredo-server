package domain

import (
	"errors"
	"fmt"
)

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

// ErrTransformWidgetFailed when a widget transformation has failed
type ErrTransformWidgetFailed struct {
	Err    error
	StdErr string
}

// NewErrTransformWidgetFailed instantiates a new ErrTransformWidgetFailed.
func NewErrTransformWidgetFailed(err error, stderr string) error {
	return &ErrTransformWidgetFailed{Err: err, StdErr: stderr}
}

func (e ErrTransformWidgetFailed) Error() string {
	return fmt.Sprintf("%s : %s", e.Err, e.StdErr)
}

// ErrDependencyNotResolved when an adapter dependency could not be resolved.
type ErrDependencyNotResolved struct {
	Adapter *Adapter
	Err     error
	Cmd     string
}

// NewErrDependencyNotResolved instantiates a new ErrDependencyNotResolved.
func NewErrDependencyNotResolved(adapter *Adapter, dependency string, err error) error {
	return &ErrDependencyNotResolved{Err: err, Cmd: dependency, Adapter: adapter}
}

func (e ErrDependencyNotResolved) Error() string {
	return fmt.Sprintf("Adapter \"%s\" dependency \"%s\" could not be resolved : %s", e.Adapter.Name, e.Cmd, e.Err)
}

// ErrParseCommandFailed when an adapter command could not be parsed.
type ErrParseCommandFailed struct {
	Adapter *Adapter
	Err     error
	Cmd     string
}

// NewErrParseCommandFailed instantiates a new ErrParseCommandFailed.
func NewErrParseCommandFailed(adapter *Adapter, cmd string, err error) error {
	return &ErrParseCommandFailed{Err: err, Cmd: cmd, Adapter: adapter}
}

func (e ErrParseCommandFailed) Error() string {
	return fmt.Sprintf("Adapter \"%s\" command \"%s\" could not be parsed : %s", e.Adapter.Name, e.Cmd, e.Err)
}
