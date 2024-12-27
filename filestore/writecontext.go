package filestore

import (
	"io"

	"github.com/metux/go-nebulon/wire"
)

type fileWriteContext struct {
	fs FileStore
	graph wire.BlockRefList
}

// 2do: uniq-sort the list
// write out the graph unencrypted
func (ctx * fileWriteContext) AddGraphRef(ref wire.BlockRef) {
	ctx.graph.AddRef(ref)
}

func (ctx * fileWriteContext) AddGraphRefs(refs [] * wire.BlockRef) {
	for _, ent := range refs {
		ctx.AddGraphRef(*ent)
	}
}

func (ctx * fileWriteContext) storeFileData(r io.Reader) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}
	buf := make([]byte, BlockSize)
	for {
		readTotal, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
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
		ctx.AddGraphRefs(reflist.Refs)
		return ctx.fs.writeBlockRefList(reflist)
	}

	new_reflist := wire.BlockRefList{}
	for reflist.Count() > BlockListMax {
		sub := wire.BlockRefList{Refs: reflist.Refs[:BlockListMax]}
		subref, err := ctx.fs.writeBlockRefList(sub)
		ctx.AddGraphRefs(sub.Refs)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		ctx.AddGraphRef(subref)
		reflist.Refs = reflist.Refs[BlockListMax:]
	}

	if reflist.Count() > 0 {
		ctx.AddGraphRefs(reflist.Refs)
		subref, err := ctx.fs.writeBlockRefList(reflist)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		ctx.graph.AddRef(subref)
	}

	return ctx.storeRefLists(new_reflist)
}
