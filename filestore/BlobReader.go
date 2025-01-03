package filestore

import (
	"bytes"
	"io"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

// io.Reader implementation for reading a Blob object as a stream
// the object is loaded lazily (only when actually read)
//
// 2do: implement background loading
type BlobReader struct {
	Ref        wire.BlockRef
	BlockStore base.BlockStore
	reader     io.Reader
}

func (r *BlobReader) Read(p []byte) (int, error) {
	if r.reader == nil {
		data, err := blockcrypt.BlockLoadDecrypt(r.BlockStore, r.Ref)
		if err != nil {
			return 0, err
		}
		r.reader = bytes.NewReader(data)
	}
	return r.reader.Read(p)
}

func (r *BlobReader) Close() error {
	r.reader = nil
	r.Ref = wire.BlockRef{}
	r.BlockStore = nil
	return nil
}

func NewBlobReader(bs base.BlockStore, ref wire.BlockRef) *BlobReader {
	reader := BlobReader{
		BlockStore: bs,
		Ref:        ref,
	}
	return &reader
}
