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

func (fs FileStore) mkWriter() FileWriteContext {
	return FileWriteContext{
		BlockStore:    fs.BlockStore,
		Cipher:        fs.Encryption,
		DataBlockSize: fs.BlockSize,
		BlockListMax:  fs.BlockListMax,
	}
}

func (fs FileStore) StoreStream(r io.Reader, header wire.Header) (wire.BlockRef, error) {
	if TraceWrite {
		defer util.TimeTrack(time.Now(), "StoreStream")
	}
	wr := fs.mkWriter()
	return wr.StoreStream(r, header)
}

func (fs FileStore) ReadStream(ref wire.BlockRef) (io.Reader, wire.Header, error) {
	reader := &fileReader{
		fs: fs,
	}

	return reader.ReadStream(ref)
}

func (fs FileStore) StoreDirectory(entries wire.BlockRefList) (wire.BlockRef, error) {
	ctx := fs.mkWriter()

	content_ref, err := ctx.StoreRefLists(entries, ctx.Cipher)
	if err != nil {
		return content_ref, err
	}

	return ctx.StoreFileControl(
		wire.FileControl{
			Content:   &content_ref,
			Directory: true})
}
