package blockstore

import (
	"crypto/sha256"
	"fmt"
)

type Key struct {
	Data []byte
}

func (k Key) String() string {
	return fmt.Sprintf("%X", k.Data)
}

func KeyForData(data []byte) Key {
	sha := sha256.Sum256(data)
	return Key{Data: sha[:]}
}
