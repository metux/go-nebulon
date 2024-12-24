package base

import (
	"github.com/metux/go-nebulon/wire"
)

type BlockStore interface {
	StoreBlock([]byte) (wire.BlockRef, error)
	LoadBlock(wire.BlockRef) ([]byte, error)
}
