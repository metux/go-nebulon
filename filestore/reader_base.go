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

func (r *ReaderBase) loadFileControl(ref wire.BlockRef) (wire.FileControl, error) {
	fctrl := wire.FileControl{}

	filehead_bin, err := r.loadBlock(ref)
	if err != nil {
		return fctrl, fmt.Errorf("failed loading FileHead [%w]", err)
	}
	filehead := wire.FileHead{}
	filehead.Unmarshal(filehead_bin)

	fctrl_bin, err := blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, filehead.Private)
	if err != nil {
		return fctrl, fmt.Errorf("error decrypting FileControl [%w]", err)
	}

	fctrl.Unmarshal(fctrl_bin)

	return fctrl, nil
}
