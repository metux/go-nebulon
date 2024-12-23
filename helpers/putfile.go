package helpers

import (
	"os"
	"time"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/util"
)

var (
	TraceStoreFile = false
)

func PutFile(fs base.IFileStore, name string, header wire.Header, fn string) (base.BlockRef, error) {
	if TraceStoreFile {
		util.TimeTrack(time.Now(), "PutFile ("+fn+")")
	}

	file, err := os.Open(fn)
	if err != nil {
		return base.BlockRef{}, err
	}
	defer file.Close()

	ref, err := fs.StoreStream(file, header)
	ref.Name = name
	return ref, err
}
