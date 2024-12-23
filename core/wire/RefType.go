package wire

import (
	"strconv"
)

func ParseRefType(s string) RefType {
	// lookup by name
	if val, ok := RefType_value[s]; ok {
		return RefType(val)
	}

	// try to decode int - if failed, assume zero (None)
	val, _ := strconv.Atoi(s)
	return RefType(val)
}
