package helpers

import (
	"fmt"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

func GetFile(fs base.FileStore, fn string, ref wire.BlockRef) (wire.Header, error) {
	if !ref.IsFile() {
		return wire.Header{}, fmt.Errorf("ref is not a file")
	}

	reader, headers, err := fs.ReadStream(ref)

	if err != nil {
		return headers, err
	}

	return headers, util.CopyStreamToFile(reader, fn)
}
