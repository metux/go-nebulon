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

func (reader *fileReader) AddRef(ref wire.BlockRef) error {
	switch ref.Type {
	case wire.RefType_Blob:
		reader.AddReader(NewBlobReader(reader.BlockStore, ref))
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

func (r *fileReader) ReadStream(ref wire.BlockRef) (io.ReadCloser, wire.Header, error) {
	fctrl, err := r.loadFileControl(ref)
	if err != nil {
		return nil, nil, err
	}

	if err = r.AddRef(*fctrl.Content); err != nil {
		return nil, nil, err
	}

	return r, fctrl.Headers, nil
}
