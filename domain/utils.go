package domain

import (
	"crypto/md5"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// Validator for structs of the domain
var validate = validator.New()

// ComputeCheckSum computes the checksum of given data.
func ComputeCheckSum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

// IsComponent checks if a string is a React component (ie. it starts with a <).
func IsComponent(str string) bool {
	return str[:1] == "<"
}
