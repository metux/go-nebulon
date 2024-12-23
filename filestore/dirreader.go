package filestore

import (
	"fmt"
	"log"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/crypt"
	"github.com/metux/go-nebulon/core/wire"
)

type DirHandle struct {
	BlockStore base.IBlockStore
	refs       []base.BlockRef
}

func (dh *DirHandle) Load(ref base.BlockRef) error {
	fctrl, err := crypt.LoadFileControl(dh.BlockStore, ref)
	if err != nil {
		panic(err)
	}
	dh.addRef(*fctrl.Content)
	return nil
}

func (dh *DirHandle) addRef(ref base.BlockRef) error {
	switch ref.Type {
	case wire.RefType_Blob:
		log.Printf("DirHandle: didnt expect blob here\n")
	case wire.RefType_RefList:
		bl, err := crypt.BlockListLoadDecrypt(dh.BlockStore, ref)
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

func (dh DirHandle) Entries() []base.BlockRef {
	return dh.refs
}
