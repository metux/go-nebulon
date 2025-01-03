package filestore

import (
	"fmt"
	"log"

	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type DirHandle struct {
	readerBase
	refs []wire.BlockRef
}

func (dh *DirHandle) Load(ref wire.BlockRef) error {
	fctrl, err := blockcrypt.LoadFileControl(dh.BlockStore, ref)
	if err != nil {
		panic(err)
	}
	dh.addRef(*fctrl.Content)
	return nil
}

func (dh *DirHandle) addRef(ref wire.BlockRef) error {
	switch ref.Type {
	case wire.RefType_Blob:
		log.Printf("DirHandle: didnt expect blob here\n")
	case wire.RefType_RefList:
		bl, err := blockcrypt.BlockListLoadDecrypt(dh.BlockStore, ref)
		if err != nil {
			return err
		}
		for _, walk := range bl.Refs {
			dh.addRef(*walk)
		}
	case wire.RefType_File:
		dh.refs = append(dh.refs, ref)
	case wire.RefType_Directory:
		dh.refs = append(dh.refs, ref)
	default:
		return fmt.Errorf("unsupported ref type %+v\n", ref.Type)
	}

	return nil
}

func (dh DirHandle) Entries() []wire.BlockRef {
	return dh.refs
}
