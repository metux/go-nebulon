package base

import (
	"github.com/metux/go-nebulon/wire"
	"io"
)

type FileStore interface {
	StoreStream(r io.Reader, headers map[string]string) (wire.BlockRef, error)
	ReadStream(wire.BlockRef) (io.Reader, map[string]string, error)
}
