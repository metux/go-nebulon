package filestore

import (
	"io"
	"log"

	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type fileWriteContext struct {
	fs    FileStore
	graph wire.BlockRefList
}

// write out the graph unencrypted
func (ctx *fileWriteContext) AddGraphRef(ref wire.BlockRef) {
	ctx.graph.AddRef(wire.BlockRef{Oid: ref.Oid})
}

func (ctx *fileWriteContext) storeFileData(r io.Reader) (wire.BlockRefList, error) {
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
	}
	return reflist, nil
}

func (ctx *fileWriteContext) storeRefLists(reflist wire.BlockRefList) (wire.BlockRef, error) {
	// store all our refs into the global graph
	for _, ent := range reflist.Refs {
		ctx.AddGraphRef(*ent)
	}

	if reflist.Count() <= BlockListMax {
		subref, err := ctx.writeBlockRefList(reflist)
		// store newly created ref into the global graph
		ctx.AddGraphRef(subref)
		return subref, err
	}

	new_reflist := wire.BlockRefList{}
	for reflist.Count() > BlockListMax {
		sub := wire.BlockRefList{Refs: reflist.Refs[:BlockListMax]}
		subref, err := ctx.writeBlockRefList(sub)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		reflist.Refs = reflist.Refs[BlockListMax:]
	}

	if reflist.Count() > 0 {
		subref, err := ctx.writeBlockRefList(reflist)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
	}

	return ctx.storeRefLists(new_reflist)
}

func (ctx *fileWriteContext) writeBlockRefList(reflist wire.BlockRefList) (wire.BlockRef, error) {
	data, err := reflist.Marshal()

	if err != nil {
		log.Println("marshal error: ", err)
		return wire.BlockRef{}, err
	}

	key, encrypted, err := blockcrypt.BlockEncrypt(ctx.fs.encryption, data)
	if err != nil {
		log.Printf("writeBlockRefList: BlockEncrypt() error %s\n", err)
		return wire.BlockRef{}, err
	}

	ref, err := ctx.fs.BlockStore.StoreBlock(encrypted)
	if err != nil {
		log.Println("error storing reflist block", err)
		return ref, err
	}

	ref.Type = wire.RefType_RefList
	ref.Cipher = ctx.fs.encryption
	ref.Key = key
	return ref, nil
}
