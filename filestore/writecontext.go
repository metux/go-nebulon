package filestore

import (
	"fmt"
	"io"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/crypt"
	"github.com/metux/go-nebulon/core/wire"
)

type FileWriteContext struct {
	// the underlying BlockStore to write into
	BlockStore base.IBlockStore
	// publicly visible block references: only used for garbage collect,
	// not suitable for reconstructing the original object, just listing
	// all other Blocks that are needed somehow
	Grabs base.BlockRefList
	// the CipherType to use for encrypting data blocks
	Cipher wire.CipherType
	// size of data blocks when storing a stream
	DataBlockSize int
	// max number of BlockRef entries when creating BlockList's
	// bigger lists will be splitted into separate BlockList objects.
	BlockListMax int
	// shall we record the block size in refs (for optimized streaming) ?
	RecordBlockSize bool
}

func (ctx *FileWriteContext) AddGrabRef(ref base.BlockRef) {
	ctx.Grabs.AddRef(ref.ToGrab())
}

func (ctx *FileWriteContext) storeFileData(r io.Reader, cipher wire.CipherType) (base.BlockRefList, uint64, error) {
	size := uint64(0)
	reflist := base.BlockRefList{}
	buf := make([]byte, ctx.DataBlockSize)
	for {
		readTotal, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				return reflist, size, err
			}
			break
		}
		size = size + uint64(readTotal)
		ref, err := ctx.writeDataBlock(buf[:readTotal], cipher)
		if err != nil {
			return reflist, size, err
		}
		reflist.AddRef(ref)
	}
	return reflist, size, nil
}

// Recursively store BlockRef list
//
// If the list is too long for one block, it's splitted into several ones, and
// their refs in turn are stored in another block - that goes on recursively
// (forming a ref tree), until we've got only one BlockRef left.
//
// The created BlockRef's are of type `RefList`. Input BlockRef's may be of
// any type.
func (ctx *FileWriteContext) StoreRefLists(reflist base.BlockRefList, cipher wire.CipherType) (base.BlockRef, error) {
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

	new_reflist := base.BlockRefList{}
	for reflist.Count() > ctx.BlockListMax {
		sub := base.BlockRefList{Refs: reflist.Refs[:ctx.BlockListMax]}
		subref, err := ctx.writeBlockRefList(sub, cipher)
		if err != nil {
			return base.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		reflist.Refs = reflist.Refs[ctx.BlockListMax:]
	}

	if reflist.Count() > 0 {
		subref, err := ctx.writeBlockRefList(reflist, cipher)
		if err != nil {
			return base.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
	}

	return ctx.StoreRefLists(new_reflist, cipher)
}

func (ctx *FileWriteContext) writeBlockRefList(reflist base.BlockRefList, cipher wire.CipherType) (base.BlockRef, error) {
	data, err := reflist.Marshal()

	if err != nil {
		return base.BlockRef{}, fmt.Errorf("reflist marshal error [%w]", err)
	}

	key, encrypted, cipher, err := crypt.BlockEncrypt(cipher, data)
	if err != nil {
		return base.BlockRef{}, fmt.Errorf("writeBlockRefList: BlockEncrypt() error [%w]", err)
	}

	ref, err := ctx.BlockStore.PutBlock(encrypted, wire.RefType_RefList)
	if err != nil {
		return ref, fmt.Errorf("writeBlockRefList: PutBlock() error [%w]", err)
	}

	ref.Type = wire.RefType_RefList
	ref.Cipher = cipher
	ref.Key = key
	return ref, nil
}

func (ctx *FileWriteContext) storeFileStream(r io.Reader, cipher wire.CipherType) (base.BlockRef, uint64, error) {
	reflist, size, err := ctx.storeFileData(r, cipher)
	if err != nil {
		return base.BlockRef{}, 0, err
	}

	content_ref, err := ctx.StoreRefLists(reflist, ctx.Cipher)
	if err != nil {
		return base.BlockRef{}, 0, err
	}

	return content_ref, size, nil
}

func (ctx *FileWriteContext) writeGrabs() (base.BlockRef, error) {
	ctx.Grabs.Sort()
	grab_ref, err := ctx.StoreRefLists(ctx.Grabs, wire.CipherType_None)
	if err != nil {
		return grab_ref, fmt.Errorf("grab write error [%w]", err)
	}
	return grab_ref, err
}

func (ctx *FileWriteContext) writeFileHead(encrypted []byte, grab_ref base.BlockRef, directory bool) (base.BlockRef, error) {
	filehead := wire.FileHead{
		Private: encrypted,
		Grabs:   &grab_ref,
	}
	filehead_bin, err := filehead.Marshal()
	if err != nil {
		return base.BlockRef{}, fmt.Errorf("error marshalling file head [%w]", err)
	}

	t := wire.RefType_File
	if directory {
		t = wire.RefType_Directory
	}

	filehead_ref, err := ctx.BlockStore.PutBlock(filehead_bin, t)
	if err != nil {
		return base.BlockRef{}, fmt.Errorf("error storing file head in blockstore [%w]", err)
	}

	return filehead_ref, nil
}

func (ctx *FileWriteContext) writeDataBlock(data []byte, cipher wire.CipherType) (base.BlockRef, error) {
	key, encrypted, cipher, err := crypt.BlockEncrypt(cipher, data)
	if err != nil {
		return base.BlockRef{}, fmt.Errorf("writeDataBlock: BlockEncrypt() error [%w]", err)
	}

	ref, err := ctx.BlockStore.PutBlock(encrypted, wire.RefType_Blob)
	if err != nil {
		return base.BlockRef{}, fmt.Errorf("writeDataBlock: storing block failed [%w]", err)
	}

	ref.Key = key
	ref.Cipher = cipher

	if ctx.RecordBlockSize {
		ref.Limit = int32(len(data))
	}

	return ref, nil
}

func (ctx *FileWriteContext) StoreStream(r io.Reader, header wire.Header) (base.BlockRef, error) {
	content_ref, size, err := ctx.storeFileStream(r, ctx.Cipher)
	if err != nil {
		return base.BlockRef{}, err
	}

	return ctx.StoreFileControl(
		wire.FileControl{
			Content:   &content_ref,
			Header:    header,
			Directory: false,
			Size:      size})
}

func (ctx *FileWriteContext) StoreFileControl(fctrl wire.FileControl) (base.BlockRef, error) {

	key, encrypted, cipher, err := crypt.EncryptFileControl(ctx.Cipher, fctrl)

	if err != nil {
		return base.BlockRef{}, err
	}

	grab_ref, err := ctx.writeGrabs()
	if err != nil {
		return base.BlockRef{}, err
	}

	filehead_ref, err := ctx.writeFileHead(encrypted, grab_ref, fctrl.Directory)
	if err != nil {
		return base.BlockRef{}, fmt.Errorf("error storing file head in blockstore [%w]", err)
	}

	filehead_ref.Cipher = cipher
	filehead_ref.Key = key
	if fctrl.Directory {
		filehead_ref.Type = wire.RefType_Directory
	} else {
		filehead_ref.Type = wire.RefType_File
	}
	return filehead_ref, nil
}
