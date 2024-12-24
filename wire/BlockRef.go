package wire

import (
	"crypto/sha256"
	"fmt"
)

func RefForBlock(data []byte) BlockRef {
	d := sha256.Sum256(data)
	return BlockRef{
		Type: RefType_Blob,
		Oid:  d[:],
	}
}

func (ref BlockRef) HexKey() string {
	return fmt.Sprintf("%s-%X", ref.Type, ref.Oid)
}

func (ref BlockRef) OID() string {
	return fmt.Sprintf("%X", ref.Oid)
}

func (ref BlockRef) Dump() string {
	return ref.HexKey()
}
