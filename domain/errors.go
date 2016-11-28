package domain

import "fmt"

import "bytes"

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
	Code    string
}

func newError(code string, message string, innerError error) error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     innerError,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s\n\t%s", e.Message, e.Err)
}

// MultipleErrors is used when multiple errors should be returned such as validation ones.
type MultipleErrors []error

func (e MultipleErrors) Error() string {
	var buf bytes.Buffer

	for _, v := range e {
		fmt.Fprintf(&buf, "%s\n", v)
	}

	return buf.String()
}

// AdapterError represents error related to an adapter.
type AdapterError struct {
	Adapter Adapter
	Err     error
}

func newAdapterError(Adapter Adapter, err error) error {
	return &AdapterError{
		Adapter: Adapter,
		Err:     err,
	}
}

func (e *AdapterError) Error() string {
	return fmt.Sprintf("%s: %s", e.Adapter.Name, e.Err)
}
