package wire

import (
	"strconv"
)

func ParseCipherType(s string) CipherType {
	// lookup by name
	if val, ok := CipherType_value[s]; ok {
		return CipherType(val)
	}

	// try to decode int - if failed, assume zero (None)
	val, _ := strconv.Atoi(s)
	return CipherType(val)
}
