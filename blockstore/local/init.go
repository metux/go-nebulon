package local

import (
	"github.com/metux/go-nebulon/blockstore/common"
	"github.com/metux/go-nebulon/core/base"
)

func init() {
	common.RegisterStoreType(StoreType,
		func(config base.BlockStoreConfig, links map[string]base.IBlockStore) (base.IBlockStore, error) {
			return NewByConfig(config, links)
		})
}
