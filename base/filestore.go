package base

import (
	"io"

	"github.com/metux/go-nebulon/wire"
)

type FileStore interface {
	StoreStream(r io.Reader, headers map[string]string) (wire.BlockRef, error)
	ReadStream(wire.BlockRef) (io.Reader, map[string]string, error)
}
