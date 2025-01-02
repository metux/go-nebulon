package helpers

import (
	"os"
	"time"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

var (
	TraceStoreFile = false
)

func PutFile(fs base.FileStore, name string, header wire.Header, fn string) (wire.BlockRef, error) {
	if TraceStoreFile {
		util.TimeTrack(time.Now(), "PutFile ("+fn+")")
	}

	file, err := os.Open(fn)
	if err != nil {
		return wire.BlockRef{}, err
	}
	defer file.Close()

	ref, err := fs.StoreStream(file, header)
	ref.Name = name
	return ref, err
}
