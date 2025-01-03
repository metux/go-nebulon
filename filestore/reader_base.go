package filestore

import (
	"github.com/metux/go-nebulon/base"
)

type readerBase struct {
	// the underlying BlockStore to write into
	BlockStore base.BlockStore
}
