package helpers

import (
	"fmt"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/util"
)

func GetFile(fs base.IFileStore, fn string, ref base.BlockRef) (wire.Header, error) {
	if !ref.IsFile() {
		return wire.Header{}, fmt.Errorf("ref is not a file")
	}

	reader, headers, _, err := fs.ReadStream(ref, 0)

	if err != nil {
		return headers, err
	}

	return headers, util.CopyStreamToFile(reader, fn)
}
