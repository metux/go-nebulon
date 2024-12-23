package util

import (
	"crypto/sha256"
)

func ContentKey(data []byte) []byte {
	s := sha256.Sum256(data)
	return s[:]
}
