package helpers

import (
	"log"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

func GetFile(fs base.FileStore, fn string, ref wire.BlockRef) (map[string]string, error) {
	reader, headers, err := fs.ReadStream(ref)
	log.Printf("Headers: %s\n", headers)

	if err != nil {
		return headers, err
	}

	return headers, util.CopyStreamToFile(reader, fn)
}
