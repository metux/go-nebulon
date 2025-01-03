package filestore

import (
	"bytes"
	"fmt"
	"io"
	"log"

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
}

func (r *BlobReader) Read(p []byte) (int, error) {
	if r.reader == nil {
		switch r.Ref.Type {
		case wire.RefType_Blob:
			data, err := blockcrypt.BlockLoadDecrypt(r.BlockStore, r.Ref)
			if err != nil {
				return 0, err
			}
			r.reader = bytes.NewReader(data)
		case wire.RefType_RefList:
			log.Printf("adding reflist ... %s\n", r.Ref.Dump())
			subreaders := make([]io.Reader, 0)
			data, err := blockcrypt.BlockLoadDecrypt(r.BlockStore, r.Ref)
			if err != nil {
				log.Printf("loading sub-ref failed: %s\n", err)
				return 0, err
			}

			// note do it in separate steps, since reflist is changed here
			bl := wire.BlockRefList{}
			err = bl.Unmarshal(data)
			if err != nil {
				log.Printf("unmarshalling sub-ref failed: %s\n", err)
				return 0, err
			}

			for _, walk := range bl.Refs {
				subreaders = append(subreaders, NewBlobReader(r.BlockStore, *walk))
			}
			r.reader = util.NewChainedReader(subreaders...)
		default:
			panic(fmt.Sprintf("unsupported ref type: %s\n", r.Ref.Type))
		}
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
