package filestore

import (
	"fmt"
	"log"

	//	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

type DirHandle struct {
	readerBase
	ents []DirEntry
	refs []wire.BlockRef
}

func (dh *DirHandle) Load(ref wire.BlockRef) error {
	log.Printf("Loading ... %s\n", ref.Dump())

	fctrl, err := dh.loadFileControl(ref)
	if err != nil {
		panic(err)
	}

	dh.addRef(*fctrl.Content)

	log.Printf("%+v\n", fctrl)

	return nil
}

func (dh *DirHandle) addRef(ref wire.BlockRef) error {
	switch ref.Type {
	case wire.RefType_Blob:
		log.Printf("DirHandle: didnt expect blob here\n")
	case wire.RefType_RefList:
		bl, err := dh.loadBlockList(ref)
		if err != nil {
			return err
		}
		for _, walk := range bl.Refs {
			dh.addRef(*walk)
		}
	case wire.RefType_File:
		dh.ents = append(dh.ents, NewDirEntry(ref))
	case wire.RefType_Directory:
		dh.ents = append(dh.ents, NewDirEntry(ref))
	default:
		return fmt.Errorf("unsupported ref type %+v\n", ref.Type)
	}

	return nil
}

func (dh *DirHandle) Entries() [] DirEntry {
	return dh.ents
}
