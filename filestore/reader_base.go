package filestore

import (
	"fmt"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type ReaderBase struct {
	// the underlying BlockStore to write into
	BlockStore base.BlockStore
}

func (r ReaderBase) loadBlockList(ref wire.BlockRef) (wire.BlockRefList, error) {
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

func (r ReaderBase) loadBlock(ref wire.BlockRef) ([]byte, error) {
	data, err := r.BlockStore.LoadBlock(ref)
	if err != nil {
		return data, err
	}

	return blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, data)
}
