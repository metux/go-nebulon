package blockcrypt

import (
	"log"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

func BlockListLoadDecrypt(bs base.BlockStore, ref wire.BlockRef) (wire.BlockRefList, error) {
	bl := wire.BlockRefList{}

	data, err := BlockLoadDecrypt(bs, ref)
	if err != nil {
		log.Printf("loading blocklist failed: %s\n", err)
		return bl, err
	}

	// note do it in separate steps, since reflist is changed here
	err = bl.Unmarshal(data)
	if err != nil {
		log.Printf("unmarshalling blocklist failed: %s\n", err)
		return bl, err
	}

	return bl, nil
}
