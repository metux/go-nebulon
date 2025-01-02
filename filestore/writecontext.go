package filestore

import (
	"fmt"
	"io"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type FileWriteContext struct {
	// the underlying BlockStore to write into
	BlockStore base.BlockStore
	// publicly visible block references: only used for garbage collect,
	// not suitable for reconstructing the original object, just listing
	// all other Blocks that are needed somehow
	Grabs wire.BlockRefList
	// the CipherType to use for encrypting data blocks
	Cipher wire.CipherType
	// size of data blocks when storing a stream
	DataBlockSize int
	// max number of BlockRef entries when creating BlockList's
	// bigger lists will be splitted into separate BlockList objects.
	BlockListMax int
}

func (ctx *FileWriteContext) AddGrabRef(ref wire.BlockRef) {
	ctx.Grabs.AddRef(wire.BlockRef{Oid: ref.Oid})
}

func (ctx *FileWriteContext) storeFileData(r io.Reader, cipher wire.CipherType) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}
	buf := make([]byte, ctx.DataBlockSize)
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

// Recursively store BlockRef list
//
// If the list is too long for one block, it's splitted into several ones, and
// their refs in turn are stored in another block - that goes on recursively
// (forming a ref tree), until we've got only one BlockRef left.
//
// The created BlockRef's are of type `RefList`. Input BlockRef's may be of
// any type.
func (ctx *FileWriteContext) StoreRefLists(reflist wire.BlockRefList, cipher wire.CipherType) (wire.BlockRef, error) {
	// store all our refs into the global grab list
	for _, ent := range reflist.Refs {
		ctx.AddGrabRef(*ent)
	}

	if reflist.Count() <= ctx.BlockListMax {
		subref, err := ctx.writeBlockRefList(reflist, cipher)
		// store newly created ref into the global grab list
		ctx.AddGrabRef(subref)
		return subref, err
	}

	new_reflist := wire.BlockRefList{}
	for reflist.Count() > ctx.BlockListMax {
		sub := wire.BlockRefList{Refs: reflist.Refs[:ctx.BlockListMax]}
		subref, err := ctx.writeBlockRefList(sub, cipher)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		reflist.Refs = reflist.Refs[ctx.BlockListMax:]
	}

	if reflist.Count() > 0 {
		subref, err := ctx.writeBlockRefList(reflist, cipher)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
	}

	return ctx.StoreRefLists(new_reflist, cipher)
}

func (ctx *FileWriteContext) writeBlockRefList(reflist wire.BlockRefList, cipher wire.CipherType) (wire.BlockRef, error) {
	data, err := reflist.Marshal()

	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("reflist marshal error [%w]", err)
	}

	key, encrypted, cipher, err := blockcrypt.BlockEncrypt(cipher, data)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("writeBlockRefList: BlockEncrypt() error [%w]", err)
	}

	ref, err := ctx.BlockStore.StoreBlock(encrypted)
	if err != nil {
		return ref, fmt.Errorf("writeBlockRefList: StoreBlock() error [%w]", err)
	}

	ref.Type = wire.RefType_RefList
	ref.Cipher = cipher
	ref.Key = key
	return ref, nil
}

func (ctx *FileWriteContext) storeFileStream(r io.Reader, cipher wire.CipherType) (wire.BlockRef, error) {
	reflist, err := ctx.storeFileData(r, cipher)
	if err != nil {
		return wire.BlockRef{}, err
	}

	content_ref, err := ctx.StoreRefLists(reflist, ctx.Cipher)
	if err != nil {
		return wire.BlockRef{}, err
	}

	return content_ref, nil
}

func (ctx *FileWriteContext) writeGrabs() (wire.BlockRef, error) {
	ctx.Grabs.Sort()
	grab_ref, err := ctx.StoreRefLists(ctx.Grabs, wire.CipherType_None)
	if err != nil {
		return grab_ref, fmt.Errorf("grab write error [%w]", err)
	}
	return grab_ref, err
}

func (ctx *FileWriteContext) writeFileHead(encrypted []byte, grab_ref wire.BlockRef) (wire.BlockRef, error) {
	filehead := wire.FileHead{
		Private: encrypted,
		Grabs:   &grab_ref,
	}
	filehead_bin, err := filehead.Marshal()
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("error marshalling file head [%w]", err)
	}

	filehead_ref, err := ctx.BlockStore.StoreBlock(filehead_bin)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("error storing file head in blockstore [%w]", err)
	}

	return filehead_ref, nil
}

func (ctx *FileWriteContext) writeDataBlock(data []byte, cipher wire.CipherType) (wire.BlockRef, error) {
	key, encrypted, cipher, err := blockcrypt.BlockEncrypt(cipher, data)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("writeDataBlock: BlockEncrypt() error [%w]", err)
	}

	ref, err := ctx.BlockStore.StoreBlock(encrypted)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("writeDataBlock: storing block failed [%w]", err)
	}

	ref.Key = key
	ref.Cipher = cipher

	return ref, nil
}

func (ctx *FileWriteContext) StoreStream(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	content_ref, err := ctx.storeFileStream(r, ctx.Cipher)
	if err != nil {
		return wire.BlockRef{}, err
	}

	return ctx.StoreFileControl(
		wire.FileControl{
			Content:   &content_ref,
			Headers:   headers,
			Directory: false})
}

func (ctx *FileWriteContext) StoreFileControl(fctrl wire.FileControl) (wire.BlockRef, error) {

	key, encrypted, cipher, err := blockcrypt.EncryptFileControl(ctx.Cipher, fctrl)

	if err != nil {
		return wire.BlockRef{}, err
	}

	grab_ref, err := ctx.writeGrabs()
	if err != nil {
		return wire.BlockRef{}, err
	}

	filehead_ref, err := ctx.writeFileHead(encrypted, grab_ref)
	if err != nil {
		return wire.BlockRef{}, fmt.Errorf("error storing file head in blockstore [%w]", err)
	}

	filehead_ref.Cipher = cipher
	filehead_ref.Key = key
	filehead_ref.Type = wire.RefType_File
	return filehead_ref, nil
}
