package helpers

import (
	"os"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

func StoreFile(fs base.FileStore, hdr map[string]string, fn string) (wire.BlockRef, error) {
	file, err := os.Open(fn)
	if err != nil {
		return wire.BlockRef{}, err
	}
	defer file.Close()

	return fs.StoreFile(file, hdr)
}
