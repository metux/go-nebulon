package filestore

import (
	"io"
	"time"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

// configuration section
var (
	TraceWrite = false
)

const (
	//	DefaultBlockSize    = 4096 * 16
	DefaultBlockSize    = 4096 * 1024
	DefaultBlockListMax = DefaultBlockSize / 80 // a blocklist entry is about 80 bytes
	DefaultCipher       = wire.CipherType_AES_CBC_ZSTD
)

type FileStore struct {
	BlockStore   base.BlockStore
	Encryption   wire.CipherType
	BlockSize    int
	BlockListMax int
}

func NewFileStore(bs base.BlockStore) base.FileStore {
	return FileStore{
		BlockStore:   bs,
		Encryption:   DefaultCipher,
		BlockSize:    DefaultBlockSize,
		BlockListMax: DefaultBlockListMax,
	}
}

func (fs FileStore) StoreStream(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	if TraceWrite {
		defer util.TimeTrack(time.Now(), "StoreStream")
	}
	context := fileWriteContext{
		fs:           fs,
		cipher:       fs.Encryption,
		blockSize:    fs.BlockSize,
		blockListMax: fs.BlockListMax,
	}
	return context.StoreStream(r, headers)
}

func (fs FileStore) ReadStream(ref wire.BlockRef) (io.Reader, map[string]string, error) {
	reader := &fileReader{
		fs: fs,
	}

	return reader.ReadStream(ref)
}
