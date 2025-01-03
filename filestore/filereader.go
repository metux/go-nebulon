package filestore

import (
	"fmt"
	"io"
	"log"
	"bytes"

	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

type fileReader struct {
	util.ChainedReader
	readerBase
}

type BlobReader struct {
	Ref wire.BlockRef
	BlockStore base.BlockStore
	Reader io.Reader
}

func (r * BlobReader) Read(p []byte) (int, error) {
	if r.Reader == nil {
		log.Printf("BlobReader: loading block")
		data, err := blockcrypt.BlockLoadDecrypt(r.BlockStore, r.Ref)
		if err != nil {
			log.Printf("BlobReader: load block error: %s\n", err)
			return 0, err
		}
		r.Reader = bytes.NewReader(data)
	}
	return r.Reader.Read(p)
}

func (reader *fileReader) AddRef(ref wire.BlockRef) error {
	_, err := reader.loadBlock(ref)
	if err != nil {
		return err
	}

	switch ref.Type {
	case wire.RefType_Blob:
//		reader.AddBytes(data)
		br := BlobReader{Ref: ref, BlockStore: reader.BlockStore}
		reader.AddReader(&br)
	case wire.RefType_RefList:
		bl, err := reader.loadBlockList(ref)
		if err != nil {
			return err
		}
		for _, walk := range bl.Refs {
			reader.AddRef(*walk)
		}
	default:
		return fmt.Errorf("unsupported ref type %+v\n", ref.Type)
	}

	return nil
}

func (r *fileReader) ReadStream(ref wire.BlockRef) (io.Reader, wire.Header, error) {
	fctrl, err := r.loadFileControl(ref)
	if err != nil {
		return nil, nil, err
	}

	if err = r.AddRef(*fctrl.Content); err != nil {
		return nil, nil, err
	}

	return r, fctrl.Headers, nil
}
