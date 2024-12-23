package base

import (
	"github.com/metux/go-nebulon/core/wire"
)

type CipherType = wire.CipherType
type BlockRef = wire.BlockRef
type BlockRefList = wire.BlockRefList

const (
	CipherType_None         = wire.CipherType_None
	CipherType_AES_CBC      = wire.CipherType_AES_CBC
	CipherType_AES_CBC_ZSTD = wire.CipherType_AES_CBC_ZSTD
)
