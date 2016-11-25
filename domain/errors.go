package domain

import "fmt"

import "bytes"

const (
	// AdapterCommandNotFound is thrown when a command could not be processed by an adapter.
	AdapterCommandNotFound = iota
	// RoomNameAlreadyExists is thrown when a room already exists with the same name.
	RoomNameAlreadyExists
	// DeviceNameAlreadyExists is thrown when a device with the same name already exists.
	DeviceNameAlreadyExists
	// DeviceConfigInvalid is thrown when a device config is invalid.
	DeviceConfigInvalid
)

// FieldError to describe errors affecting one field.
type FieldError struct {
	Resource string
	Field    string
	Code     string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s validation failed for %s with the code %s", e.Resource, e.Field, e.Code)
}

// Error of the domain
type Error struct {
	Err     error
	Message string
	Code    int
}

func newError(code int, message string, innerError error) error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     innerError,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %s", e.Code, e.Message)
}

// MultipleErrors is used when multiple errors should be returned.
type MultipleErrors []error

func (e MultipleErrors) Error() string {
	var buf bytes.Buffer

	for _, v := range e {
		fmt.Fprintf(&buf, "%s\n", v.Error())
	}

	return buf.String()
}

var (
	errAdapterCommandNotFound  = newError(AdapterCommandNotFound, "Adapter command not found", nil)
	errRoomNameAlreadyExists   = newError(RoomNameAlreadyExists, "Room name already exists", nil)
	errDeviceNameAlreadyExists = newError(DeviceNameAlreadyExists, "Device name already exists", nil)
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
