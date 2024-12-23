package base

import (
	"time"
)

type BlockInfo struct {
	Ref        BlockRef
	CreateTime time.Time
	ModTime    time.Time
	Size       int64
	Present    int64
	Fetching   bool
}
