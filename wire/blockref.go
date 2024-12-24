package wire

import (
	"fmt"
	"crypto/sha256"
)

func RefForBlock(data []byte) BlockRef {
	d := sha256.Sum256(data)
	return BlockRef{
		Type: RefType_Blob,
		Key: d[:],
	}
}

func (ref BlockRef) HexKey() string {
	return fmt.Sprintf("%X", ref.Key)
}
