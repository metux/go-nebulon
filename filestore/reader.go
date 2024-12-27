package filestore

import (
	"fmt"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

type FileReader struct {
	util.ChainedReader
	fs FileStore
}

func (reader *FileReader) AddRef(ref wire.BlockRef) error {
	data, err := reader.fs.LoadBlock(ref)
	if err != nil {
		return err
	}

	switch ref.Type {
	case wire.RefType_Blob:
		reader.AddBytes(data)
	case wire.RefType_RefList:
		bl, err := reader.fs.loadBlockList(ref)
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
