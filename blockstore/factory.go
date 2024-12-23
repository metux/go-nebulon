package blockstore

import (
	// link in the drivers so they can register themselves
	_ "github.com/metux/go-nebulon/blockstore/cache"
	_ "github.com/metux/go-nebulon/blockstore/grpc"
	_ "github.com/metux/go-nebulon/blockstore/http"
	_ "github.com/metux/go-nebulon/blockstore/local"

	"github.com/metux/go-nebulon/blockstore/common"
)

type Constructor = common.Constructor

var (
	CreateStore       = common.CreateStore
	NewStoreByConfig  = common.CreateStore
	RegisterStoreType = common.RegisterStoreType
)
