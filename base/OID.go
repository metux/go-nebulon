package base

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type OID struct {
	Data []byte
}

func (oid OID) String() string {
	return fmt.Sprintf("%X", oid.Data)
}

func (oid OID) Equals(in OID) bool {
	return bytes.Equal(oid.Data, in.Data)
}

func OIDForBlock(data []byte) OID {
	sha := sha256.Sum256(data)
	return OID{Data: sha[:]}
}
