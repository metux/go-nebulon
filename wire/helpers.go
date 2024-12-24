package wire

import (
	"fmt"
	"crypto/sha256"
)

func RefForBlock(data []byte) BlockRef {
	d := sha256.Sum256(data)
	return BlockRef{
		Type: RefType_Blob,
		Data: d[:],
	}
}

func RefToString(ref BlockRef) string {
	return fmt.Sprintf("%X", ref.Data)
}
