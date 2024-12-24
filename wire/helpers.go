package wire

import (
    "crypto/sha256"
)

func RefForBlock(data []byte) BlockRef {
	d := sha256.Sum256(data)
	return BlockRef {
		Type: RefType_Blob,
		Data: d[:],
	}
}
