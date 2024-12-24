package base

import (
	"io"
	"github.com/metux/go-nebulon/wire"
)

type FileStore interface {
	StoreFile (r io.Reader, headers map[string]string) (wire.BlockRef, error)
	ReadFile (wire.BlockRef) (io.Reader, map[string]string, error)
}
