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

func (ref BlockRef) Dump() string {
	return fmt.Sprintf("%s:%X:%X", ref.Type, ref.Oid, ref.Key)
}

func (ref BlockRef) OID() string {
	return fmt.Sprintf("%X", ref.Oid)
}
