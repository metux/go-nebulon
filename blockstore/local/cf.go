package local

import (
	"github.com/metux/go-nebulon/core/base"
)

const (
	IterateChanSize = 500
	IterateDirSize  = 500
	StoreType       = base.BlockStoreType_LocalFS
)

var (
	TraceWrite         = false
	VirtualRoot string = "."
)
