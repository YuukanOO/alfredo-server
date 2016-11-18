package domain

import (
	"errors"
	"fmt"
)

// TODO
type DomainError struct {
	Code    string
	Message string
}

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
	Adapter *Adapter
	Err     error
	StdErr  string
	Widget  string
}

func newErrTransformWidgetFailed(adapter *Adapter, widget string, err error, stderr string) error {
	return &ErrTransformWidgetFailed{
		Adapter: adapter,
		Widget:  widget,
		Err:     err,
		StdErr:  stderr,
	}
}

func (e ErrTransformWidgetFailed) Error() string {
	return fmt.Sprintf("Adapter \"%s\" widget \"%s\" could not be transformed %s : %s", e.Adapter.Name, e.Widget, e.Err, e.StdErr)
}

// ErrDependencyNotResolved when an adapter dependency could not be resolved.
type ErrDependencyNotResolved struct {
	Adapter *Adapter
	Err     error
	Cmd     string
}

func newErrDependencyNotResolved(adapter *Adapter, dependency string, err error) error {
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

func newErrParseCommandFailed(adapter *Adapter, cmd string, err error) error {
	return &ErrParseCommandFailed{Err: err, Cmd: cmd, Adapter: adapter}
}

func (e ErrParseCommandFailed) Error() string {
	return fmt.Sprintf("Adapter \"%s\" command \"%s\" could not be parsed : %s", e.Adapter.Name, e.Cmd, e.Err)
}
