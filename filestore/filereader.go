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
	fctrl, err := r.loadFileControl(ref)
	if err != nil {
		return nil, nil, err
	}

	blobreader := NewBlobReader(r.BlockStore, *fctrl.Content)
	r.AddReader(blobreader)

	return r, fctrl.Headers, nil
}
