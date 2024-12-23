package base

import (
	"io"

	"github.com/metux/go-nebulon/core/wire"
)

// FIXME: add IFileReader interface ?
type IFileStore interface {
	StoreStream(r io.Reader, header wire.Header) (BlockRef, error)
	ReadStream(BlockRef, uint64) (io.ReadCloser, wire.Header, uint64, error)
	StoreDirectory(BlockRefList) (BlockRef, error)
	ReadDirectory(BlockRef) ([]BlockRef, error)
	Grabs(BlockRef) BlockRefStream
	ReadFileControl(BlockRef) (wire.FileControl, error)
}
