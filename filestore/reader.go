package filestore

import (
	"errors"
	"fmt"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
	"log"
)

type FileReader struct {
	util.ChainedReader
	fs FileStore
}

func (reader *FileReader) AddRef(ref wire.BlockRef) error {
	data, err := reader.fs.LoadBlock(ref)
	if err != nil {
		log.Printf("readerAddRef: failed reading block: %s\n", err)
		return err
	}

	switch ref.Type {
	case wire.RefType_Blob:
		reader.AddBytes(data)
	case wire.RefType_RefList:
		bl, err := reader.fs.LoadBlockList(ref)
		if err != nil {
			log.Printf("readerAddFile: failed reading block list %s\n", err)
			return err
		}
		log.Printf("BLOCK REF LIST %+v\n", bl.Dump())

		for idx, walk := range bl.Refs {
			log.Printf("block ref ent %d -- %s\n", idx, walk.Dump())
			reader.AddRef(*walk)
		}
	default:
		log.Printf("readerAddRef: unsupported ref type %+v (default)\n", ref.Type)
		return errors.New(fmt.Sprintf("unsupported ref type %+v\n", ref.Type))
	}

	return nil
}
