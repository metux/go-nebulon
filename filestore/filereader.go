package filestore

import (
	"fmt"
	"io"

	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

type fileReader struct {
	util.ChainedReader
	readerBase
}

type BlobReader struct {
	Ref wire.BlockRef
}

func (reader *fileReader) AddRef(ref wire.BlockRef) error {
	data, err := reader.loadBlock(ref)
	if err != nil {
		return err
	}

	switch ref.Type {
	case wire.RefType_Blob:
		reader.AddBytes(data)
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
