package filestore

import (
	"bytes"
	"fmt"
	"io"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/util"
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
	header     wire.Header
}

func (r *BlobReader) Read(p []byte) (int, error) {
	if err := r.init(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

func (r *BlobReader) init() error {
	if r.reader != nil {
		return nil
	}

	switch r.Ref.Type {
	case wire.RefType_Blob:
		data, err := blockcrypt.BlockLoadDecrypt(r.BlockStore, r.Ref)
		if err != nil {
			return err
		}
		r.reader = bytes.NewReader(data)
	case wire.RefType_RefList:
		bl, err := blockcrypt.BlockListLoadDecrypt(r.BlockStore, r.Ref)
		if err != nil {
			return err
		}
		subreaders := make([]io.Reader, 0)
		for _, walk := range bl.Refs {
			subreaders = append(subreaders, NewBlobReader(r.BlockStore, *walk))
		}
		r.reader = util.NewChainedReader(subreaders...)
	case wire.RefType_File:
		fctrl, err := blockcrypt.LoadFileControl(r.BlockStore, r.Ref)
		if err != nil {
			return err
		}
		r.reader = NewBlobReader(r.BlockStore, *fctrl.Content)
		r.header = fctrl.Header
	default:
		return fmt.Errorf("unsupported ref type: %s\n", r.Ref.Type)
	}
	return nil
}

func (r *BlobReader) Close() error {
	r.reader = nil
	r.Ref = wire.BlockRef{}
	r.BlockStore = nil
	return nil
}

func (r *BlobReader) GetHeader() (wire.Header, error) {
	if err := r.init(); err != nil {
		return wire.Header{}, err
	}
	return r.header, nil
}

func NewBlobReader(bs base.BlockStore, ref wire.BlockRef) *BlobReader {
	reader := BlobReader{
		BlockStore: bs,
		Ref:        ref,
	}
	return &reader
}
