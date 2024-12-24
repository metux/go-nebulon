package base

import (
	"crypto/sha256"
	"fmt"
)

type OID struct {
	Data []byte
}

func (oid OID) String() string {
	return fmt.Sprintf("%X", oid.Data)
}

func OIDForBlock(data []byte) OID {
	sha := sha256.Sum256(data)
	return OID{Data: sha[:]}
}
