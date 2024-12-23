package helpers

import (
	"os"
	"github.com/metux/go-nebulon/base"
)

func StoreFile(fs base.FileStore, hdr map[string]string, fn string) (base.OID, error) {
	file, err := os.Open(fn)
	if err != nil {
		return base.OID{}, err
	}
	defer file.Close()

	return fs.StoreFile(file, hdr)
}
