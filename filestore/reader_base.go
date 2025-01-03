package filestore

import (
	"fmt"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type readerBase struct {
	// the underlying BlockStore to write into
	BlockStore base.BlockStore
}

func (r readerBase) loadBlockList(ref wire.BlockRef) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}

	encrypted, err := r.BlockStore.LoadBlock(ref)
	if err != nil {
		return reflist, fmt.Errorf("failed loading blocklist block [%w]", err)
	}

	data, err := blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, encrypted)
	if err != nil {
		return reflist, fmt.Errorf("failed decrypting blocklist [%w]", err)
	}

	// note do it in separate steps, since reflist is changed here
	err = reflist.Unmarshal(data)
	return reflist, err
}

func (r readerBase) loadBlock(ref wire.BlockRef) ([]byte, error) {
	return blockcrypt.BlockLoadDecrypt(r.BlockStore, ref)
}

func (r *readerBase) loadFileControl(ref wire.BlockRef) (wire.FileControl, error) {
	// load the index block -- strip off the, that's later used used to decrypt the FileControl block
	filehead_ref := wire.BlockRef{
		Oid:  ref.Oid,
		Type: ref.Type,
	}

	filehead_bin, err := r.loadBlock(filehead_ref)
	if err != nil {
		return wire.FileControl{}, fmt.Errorf("failed loading FileHead [%w]", err)
	}

	filehead, err := wire.FileHeadUnmarshal(filehead_bin)
	if err != nil {
		return wire.FileControl{}, fmt.Errorf("failed unmarshalling FileHead [%w]", err)
	}

	return blockcrypt.DecryptFileControl(ref.Cipher, ref.Key, filehead.Private)
}
