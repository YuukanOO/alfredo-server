package domain

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"sort"
)

// ComputeCheckSum computes the checksum of given data.
func ComputeCheckSum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

// IsComponent checks if a string is a React component (ie. it starts with a <).
func IsComponent(str string) bool {
	return str[:1] == "<"
}

// SortedKeys returns a sorted array of keys of the map.
func SortedKeys(val map[string]string) []string {
	keys := make([]string, len(val))
	i := 0

	for k := range val {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	return keys
}

// SortedValues retrieve values of a map sorted by the given array of keys.
func SortedValues(sortedKeys []string, val map[string]string) []string {
	values := make([]string, len(sortedKeys))

	for i, k := range sortedKeys {
		values[i] = k + ":" + val[k]
	}

	return values
}

// GetBytes convert an arbitrary interface to a bytes representation.
func GetBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
