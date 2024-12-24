package base

import (
	"io"
	"github.com/metux/go-nebulon/wire"
)

type FileStore interface {
    StoreFile (r io.Reader, headers map[string]string) (wire.BlockRef, error)
    ReadFile (oid OID) (io.Reader, map[string]string, error)
}
