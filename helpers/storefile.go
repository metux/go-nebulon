package helpers

import (
	"os"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

func StoreFile(fs base.FileStore, name string, hdr map[string]string, fn string) (wire.BlockRef, error) {
	file, err := os.Open(fn)
	if err != nil {
		return wire.BlockRef{}, err
	}
	defer file.Close()

	ref, err := fs.StoreStream(file, hdr)
	ref.Name = name
	return ref, err
}
