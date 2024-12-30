package filestore

import (
	"fmt"
	"io"

	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type fileWriteContext struct {
	fs           FileStore
	graph        wire.BlockRefList
	cipher       wire.CipherType
	blockSize    int
	blockListMax int
}

func (ctx *fileWriteContext) AddGraphRef(ref wire.BlockRef) {
	ctx.graph.AddRef(wire.BlockRef{Oid: ref.Oid})
}

func (ctx *fileWriteContext) storeFileData(r io.Reader, cipher wire.CipherType) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}
	buf := make([]byte, ctx.blockSize)
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

	if reflist.Count() <= ctx.blockListMax {
		subref, err := ctx.writeBlockRefList(reflist, cipher)
		// store newly created ref into the global graph
		ctx.AddGraphRef(subref)
		return subref, err
	}

	new_reflist := wire.BlockRefList{}
	for reflist.Count() > ctx.blockListMax {
		sub := wire.BlockRefList{Refs: reflist.Refs[:ctx.blockListMax]}
		subref, err := ctx.writeBlockRefList(sub, cipher)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		reflist.Refs = reflist.Refs[ctx.blockListMax:]
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
		return wire.BlockRef{}, fmt.Errorf("reflist marshal error [%w]", err)
	}

	key, encrypted, err := blockcrypt.BlockEncrypt(cipher, data)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("writeBlockRefList: BlockEncrypt() error [%w]", err)
	}

	ref, err := ctx.fs.BlockStore.StoreBlock(encrypted)
	if err != nil {
		return ref, fmt.Errorf("writeBlockRefList: StoreBlock() error [%w]", err)
	}

	ref.Type = wire.RefType_RefList
	ref.Cipher = ctx.cipher
	ref.Key = key
	return ref, nil
}

func (ctx *fileWriteContext) storeFileStream(r io.Reader, cipher wire.CipherType) (wire.BlockRef, error) {
	reflist, err := ctx.storeFileData(r, cipher)
	if err != nil {
		return wire.BlockRef{}, err
	}

	content_ref, err := ctx.storeRefLists(reflist, ctx.cipher)
	if err != nil {
		return wire.BlockRef{}, err
	}

	return content_ref, nil
}

func (ctx *fileWriteContext) writeGraph() (wire.BlockRef, error) {
	ctx.graph.Sort()
	graph_ref, err := ctx.storeRefLists(ctx.graph, wire.CipherType_None)
	if err != nil {
		return graph_ref, fmt.Errorf("graph write error [%w]", err)
	}
	return graph_ref, err
}

func (ctx *fileWriteContext) writeFileHead(encrypted []byte, graph_ref wire.BlockRef) (wire.BlockRef, error) {
	filehead := wire.FileHead{
		Private: encrypted,
		Graph:   &graph_ref,
	}
	filehead_bin, err := filehead.Marshal()
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("error marshalling file head [%w]", err)
	}

	filehead_ref, err := ctx.fs.BlockStore.StoreBlock(filehead_bin)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("error storing file head in blockstore [%w]", err)
	}

	return filehead_ref, nil
}

func (ctx *fileWriteContext) writeDataBlock(data []byte, cipher wire.CipherType) (wire.BlockRef, error) {
	key, encrypted, err := blockcrypt.BlockEncrypt(cipher, data)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("writeDataBlock: BlockEncrypt() error [%w]", err)
	}

	ref, err := ctx.fs.BlockStore.StoreBlock(encrypted)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("writeDataBlock: storing block failed [%w]", err)
	}

	ref.Key = key
	ref.Cipher = cipher

	return ref, nil
}

func (ctx *fileWriteContext) StoreStream(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	content_ref, err := ctx.storeFileStream(r, ctx.cipher)

	key, encrypted, err := blockcrypt.EncryptFileControl(content_ref, headers, ctx.cipher)

	if err != nil {
		return wire.BlockRef{}, err
	}

	graph_ref, err := ctx.writeGraph()
	if err != nil {
		return wire.BlockRef{}, err
	}

	filehead_ref, err := ctx.writeFileHead(encrypted, graph_ref)
	if err != nil {
		return content_ref, fmt.Errorf("error storing file head in blockstore [%w]", err)
	}

	filehead_ref.Cipher = ctx.cipher
	filehead_ref.Key = key
	filehead_ref.Type = wire.RefType_File
	return filehead_ref, nil
}
