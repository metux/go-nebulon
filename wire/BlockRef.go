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
	s := fmt.Sprintf("%s:%X:%X", ref.Type, ref.Oid, ref.Key)
	if ref.Name != "" {
		s = s + " <" + ref.Name + ">"
	}
	return s
}

func (ref BlockRef) OID() string {
	return fmt.Sprintf("%X", ref.Oid)
}

func (ref BlockRef) IsDir() bool {
	return ref.Type == RefType_Directory
}

func (ref BlockRef) IsFile() bool {
	return ref.Type == RefType_File
}
