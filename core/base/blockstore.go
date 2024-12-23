package base

import (
	"github.com/metux/go-nebulon/core/wire"
)

type IBlockStore interface {
	PutBlock([]byte, wire.RefType) (BlockRef, error)
	GetBlock(BlockRef) ([]byte, error)
	DeleteBlock(BlockRef) error
	KeepBlock(BlockRef) error
	IterateBlocks() BlockRefStream
	PeekBlock(ref BlockRef, fetch int) (BlockInfo, error)
	Ping() error
}
