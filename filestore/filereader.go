package filestore

import (
	"io"

	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

type fileReader struct {
	util.ChainedReader
	readerBase
}

func (r *fileReader) ReadStream(ref wire.BlockRef) (io.ReadCloser, wire.Header, error) {
	blobreader := NewBlobReader(r.BlockStore, ref)
	r.AddReader(blobreader)

	h, err := blobreader.GetHeader()
	if err != nil {
		return r, h, err
	}

	return r, h, nil
}
