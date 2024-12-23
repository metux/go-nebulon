package filestore

import (
	"bytes"
	"fmt"
	"io"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/crypt"
	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/util"
)

// io.Reader implementation for reading a Blob object as a stream
// the object is loaded lazily (only when actually read)
//
// 2do: implement background loading
type BlobReader struct {
	Ref        base.BlockRef
	BlockStore base.IBlockStore
	reader     io.Reader
	fctrl      wire.FileControl
}

func (r *BlobReader) Read(p []byte) (int, error) {
	if err := r.init(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

// FIXME: support offset and limit
func (r *BlobReader) init() error {
	if r.reader != nil {
		return nil
	}

	switch r.Ref.Type {
	case wire.RefType_Blob:
		data, err := crypt.BlockLoadDecrypt(r.BlockStore, r.Ref)
		if err != nil {
			return err
		}
		r.reader = bytes.NewReader(data)
	case wire.RefType_RefList:
		bl, err := crypt.BlockListLoadDecrypt(r.BlockStore, r.Ref)
		if err != nil {
			return err
		}
		subreaders := make([]io.Reader, 0)
		for _, walk := range bl.Refs {
			subreaders = append(subreaders, NewBlobReader(r.BlockStore, *walk))
		}
		r.reader = util.NewChainedReader(subreaders...)
	case wire.RefType_File:
		fctrl, err := crypt.LoadFileControl(r.BlockStore, r.Ref)
		if err != nil {
			return err
		}
		r.reader = NewBlobReader(r.BlockStore, *fctrl.Content)
		r.fctrl = fctrl
	default:
		return fmt.Errorf("unsupported ref type: %s\n", r.Ref.Type)
	}
	return nil
}

func (r *BlobReader) Close() error {
	r.reader = nil
	r.Ref = base.BlockRef{}
	r.BlockStore = nil
	r.fctrl = wire.FileControl{}
	return nil
}

func (r *BlobReader) GetHeader() (wire.Header, uint64, error) {
	if err := r.init(); err != nil {
		return wire.Header{}, 0, err
	}
	return r.fctrl.Header, r.fctrl.Size, nil
}

func (r *BlobReader) GetFileControl() (wire.FileControl, error) {
	if err := r.init(); err != nil {
		return wire.FileControl{}, err
	}
	return r.fctrl, nil
}

func NewBlobReader(bs base.IBlockStore, ref base.BlockRef) *BlobReader {
	reader := BlobReader{
		BlockStore: bs,
		Ref:        ref,
	}
	return &reader
}
