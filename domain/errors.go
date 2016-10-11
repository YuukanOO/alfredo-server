package domain

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
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

// ErrCommandFailed when a command has failed.
type ErrCommandFailed struct {
	Cmd    *exec.Cmd
	Err    error
	StdErr string
}

// NewErrCommandFailed instantiates a new ErrCommandFailed.
func NewErrCommandFailed(cmd *exec.Cmd, err error, stderr string) error {
	return ErrCommandFailed{Cmd: cmd, Err: err, StdErr: stderr}
}

func (e ErrCommandFailed) Error() string {
	return fmt.Sprintf("%s %s\n%s\n%s", e.Cmd.Path, strings.Join(e.Cmd.Args, " "), e.Err, e.StdErr)
}
