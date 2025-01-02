package filestore

import (
	"log"

	//	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

type DirHandle struct {
	ReaderBase
	Ref wire.BlockRef
}

func (dh *DirHandle) Load(ref wire.BlockRef) error {
	// need to ignore cipher and key here
	dirhead_ref := wire.BlockRef{
		Oid:  ref.Oid,
		Type: ref.Type,
	}

	log.Printf("%+v\n", dirhead_ref)
	return nil
}
