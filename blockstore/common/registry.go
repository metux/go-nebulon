package common

import (
	"github.com/metux/go-nebulon/core/base"
)

type Constructor = func(config base.BlockStoreConfig, links map[string]base.IBlockStore) (base.IBlockStore, error)

var (
	constructors map[string]Constructor
)

func RegisterStoreType(t string, c Constructor) {
	if constructors == nil {
		constructors = map[string]Constructor{}
	}
	constructors[t] = c
}

func CreateStore(cf base.BlockStoreConfig, links map[string]base.IBlockStore) (base.IBlockStore, error) {
	if c, ok := constructors[cf.Type]; ok {
		return c(cf, links)
	}
	return nil, base.ErrUnsupportedStore
}
