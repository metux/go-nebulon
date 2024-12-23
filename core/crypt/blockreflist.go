package crypt

import (
	"github.com/metux/go-nebulon/core/base"
)

func BlockListLoadDecrypt(bs base.IBlockStore, ref base.BlockRef) (base.BlockRefList, error) {
	bl := base.BlockRefList{}

	data, err := BlockLoadDecrypt(bs, ref)
	if err != nil {
		return bl, err
	}

	// note do it in separate steps, since reflist is changed here
	err = bl.Unmarshal(data)
	if err != nil {
		return bl, err
	}

	return bl, nil
}
