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

func (ctx *fileWriteContext) storeFileData(r io.Reader, cipher wire.CipherType) (wire.BlockRefList, error) {
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
		ref, err := ctx.writeDataBlock(buf[:readTotal], cipher)
		if err != nil {
			return reflist, err
		}
		reflist.AddRef(ref)
	}
	return reflist, nil
}

func (ctx *fileWriteContext) storeRefLists(reflist wire.BlockRefList, cipher wire.CipherType) (wire.BlockRef, error) {
	// store all our refs into the global graph
	for _, ent := range reflist.Refs {
		ctx.AddGraphRef(*ent)
	}

	if reflist.Count() <= BlockListMax {
		subref, err := ctx.writeBlockRefList(reflist, cipher)
		// store newly created ref into the global graph
		ctx.AddGraphRef(subref)
		return subref, err
	}

	new_reflist := wire.BlockRefList{}
	for reflist.Count() > BlockListMax {
		sub := wire.BlockRefList{Refs: reflist.Refs[:BlockListMax]}
		subref, err := ctx.writeBlockRefList(sub, cipher)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		reflist.Refs = reflist.Refs[BlockListMax:]
	}

	if reflist.Count() > 0 {
		subref, err := ctx.writeBlockRefList(reflist, cipher)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
	}

	return ctx.storeRefLists(new_reflist, cipher)
}

func (ctx *fileWriteContext) writeBlockRefList(reflist wire.BlockRefList, cipher wire.CipherType) (wire.BlockRef, error) {
	data, err := reflist.Marshal()

	if err != nil {
		log.Println("marshal error: ", err)
		return wire.BlockRef{}, err
	}

	key, encrypted, err := blockcrypt.BlockEncrypt(cipher, data)
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

func (ctx *fileWriteContext) storeFileStream(r io.Reader, cipher wire.CipherType) (wire.BlockRef, error) {
	reflist, err := ctx.storeFileData(r, cipher)
	if err != nil {
		return wire.BlockRef{}, err
	}

	content_ref, err := ctx.storeRefLists(reflist, ctx.fs.encryption)
	if err != nil {
		log.Printf("storeRefLists() error %s\n", err)
		return wire.BlockRef{}, err
	}

	return content_ref, nil
}

func (ctx *fileWriteContext) writeGraph() (wire.BlockRef, error) {
	ctx.graph.Sort()
	graph_ref, err := ctx.storeRefLists(ctx.graph, wire.CipherType_None)
	if err != nil {
		log.Printf("Graph write error: %s\n", err)
	}
	log.Printf("Graph ref: %s\n", graph_ref.Dump())
	return graph_ref, err
}

func (ctx *fileWriteContext) writeFileHead(encrypted []byte, graph_ref wire.BlockRef) (wire.BlockRef, error) {
	filehead := wire.FileHead{
		Private: encrypted,
		Graph:   &graph_ref,
	}
	filehead_bin, err := filehead.Marshal()
	if err != nil {
		log.Printf("error marshalling file head: %s\n", err)
		return wire.BlockRef{}, err
	}

	filehead_ref, err := ctx.fs.BlockStore.StoreBlock(filehead_bin)
	if err != nil {
		log.Printf("error storing file head in blockstore %s\n", err)
		return wire.BlockRef{}, err
	}

	log.Printf("file head ref: %X\n", filehead_ref.Oid)

	return filehead_ref, nil
}

func (ctx *fileWriteContext) writeDataBlock(data []byte, cipher wire.CipherType) (wire.BlockRef, error) {
	key, encrypted, err := blockcrypt.BlockEncrypt(cipher, data)
	if err != nil {
		log.Printf("writeDataBlock: BlockEncrypt() error %s\n", err)
		return wire.BlockRef{}, err
	}

	ref, err := ctx.fs.BlockStore.StoreBlock(encrypted)
	if err != nil {
		log.Printf("writeDataBlock: storing block failed %s\n", err)
		return wire.BlockRef{}, err
	}

	ref.Key = key
	ref.Cipher = cipher

	return ref, nil
}
