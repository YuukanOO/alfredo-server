package domain

import "fmt"
import "crypto/md5"

// ComputeCheckSum computes the checksum of given data.
func ComputeCheckSum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
