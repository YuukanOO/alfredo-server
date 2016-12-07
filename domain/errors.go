package domain

import "fmt"

import "bytes"
import "gopkg.in/go-playground/validator.v9"
import "reflect"

// FieldError to describe errors affecting one field.
type FieldError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s validation failed for %s with the code %s", e.Resource, e.Field, e.Code)
}

// Error of the domain
type Error struct {
	Err     MultipleErrors `json:"errors"`
	Message string         `json:"message"`
	Code    string         `json:"code"`
}

func newError(code string, message string, innerError error) error {
	retErr := &Error{
		Code:    code,
		Message: message,
	}

	// We must always returns an instance of MultipleErrors so the format will be easier
	switch innerError.(type) {
	case MultipleErrors:
		retErr.Err = innerError.(MultipleErrors)
	default:
		retErr.Err = MultipleErrors{innerError}
	}

	return retErr
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

// ErrValidationFailed represents a validation error
const ErrValidationFailed = "ErrValidationFailed"

func newValidationErrors(resource interface{}, err error) error {
	valErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return newError(ErrValidationFailed, "Validation failed", err)
	}

	res := reflect.ValueOf(resource)

	// If pointer get the underlying element≤
	for res.Kind() == reflect.Ptr {
		res = res.Elem()
	}

	resType := res.Type()

	retErrors := make(MultipleErrors, len(valErrors))

	for i, curErr := range valErrors {
		fieldName := curErr.Field()
		field, ok := resType.FieldByName(fieldName)

		if ok {
			fieldName = field.Tag.Get("bson")
		}

		retErrors[i] = &FieldError{
			Resource: resType.Name(),
			Field:    fieldName,
			Code:     curErr.ActualTag(),
		}
	}

	return newError(ErrValidationFailed, "Validation failed", retErrors)
}
