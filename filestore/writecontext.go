package filestore

import (
	"io"
	"log"

//	"github.com/metux/go-nebulon/base"
//	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type fileWriteContext struct {
	fs FileStore
	graph wire.BlockRefList
}

func (ctx * fileWriteContext) storeFileData(r io.Reader) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}
	buf := make([]byte, BlockSize)
	for {
		readTotal, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return reflist, err
			}
			break
		}
		ref, err := ctx.fs.writeDataBlock(buf[:readTotal])
		if err != nil {
			return reflist, err
		}
		reflist.AddRef(ref)
		ctx.graph.AddRef(ref)
	}
	return reflist, nil
}

func (ctx * fileWriteContext) storeRefLists(reflist wire.BlockRefList) (wire.BlockRef, error) {
	if reflist.Count() <= BlockListMax {
		return ctx.fs.writeBlockRefList(reflist)
	}

	new_reflist := wire.BlockRefList{}
	for reflist.Count() > BlockListMax {
		sub := wire.BlockRefList{Refs: reflist.Refs[:BlockListMax]}
		subref, err := ctx.fs.writeBlockRefList(sub)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		ctx.graph.AddRef(subref)
		reflist.Refs = reflist.Refs[BlockListMax:]
	}

	if reflist.Count() > 0 {
		subref, err := ctx.fs.writeBlockRefList(reflist)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		ctx.graph.AddRef(subref)
	}

	return ctx.storeRefLists(new_reflist)
}
