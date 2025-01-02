package base

import (
	"io"

	"github.com/metux/go-nebulon/wire"
)

type FileStore interface {
	StoreStream(r io.Reader, header wire.Header) (wire.BlockRef, error)
	ReadStream(wire.BlockRef) (io.Reader, wire.Header, error)
	StoreDirectory(wire.BlockRefList) (wire.BlockRef, error)
}
