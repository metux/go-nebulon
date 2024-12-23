package filestore

import (
	"io"
	"time"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/util"
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
	BlockStore   base.IBlockStore
	Encryption   wire.CipherType
	BlockSize    int
	BlockListMax int
}

func NewFileStore(bs base.IBlockStore) FileStore {
	if bs == nil {
		panic("blockstore = nil")
	}
	return FileStore{
		BlockStore:   bs,
		Encryption:   DefaultCipher,
		BlockSize:    DefaultBlockSize,
		BlockListMax: DefaultBlockListMax,
	}
}

func (fs FileStore) mkWriter() FileWriteContext {
	return FileWriteContext{
		BlockStore:      fs.BlockStore,
		Cipher:          fs.Encryption,
		DataBlockSize:   fs.BlockSize,
		BlockListMax:    fs.BlockListMax,
		RecordBlockSize: true,
	}
}

func (fs FileStore) StoreStream(r io.Reader, header wire.Header) (base.BlockRef, error) {
	if TraceWrite {
		defer util.TimeTrack(time.Now(), "FileStore::StoreStream")
	}
	wr := fs.mkWriter()
	return wr.StoreStream(r, header)
}

// FIXME: switch to int64 ?
// FIXME: check for correct ref type
func (fs FileStore) ReadStream(ref base.BlockRef, offset uint64) (io.ReadCloser, wire.Header, uint64, error) {
	blobreader := NewBlobReader(fs.BlockStore, ref)

	hdr, size, err := blobreader.GetHeader()
	if err != nil {
		return blobreader, hdr, size, err
	}

	// FIXME: find a more effient way
	if offset != 0 {
		io.CopyN(io.Discard, blobreader, int64(offset))
	}

	return blobreader, hdr, size, nil
}

// FIXME: check for correct ref type
func (fs FileStore) ReadFileControl(ref base.BlockRef) (wire.FileControl, error) {
	blobreader := NewBlobReader(fs.BlockStore, ref)
	fctrl, err := blobreader.GetFileControl()
	blobreader.Close()
	return fctrl, err
}

func (fs FileStore) StoreDirectory(entries base.BlockRefList) (base.BlockRef, error) {
	if TraceWrite {
		defer util.TimeTrack(time.Now(), "FileStore::StoreDirectory")
	}

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

func (fs FileStore) ReadDirectory(ref base.BlockRef) ([]base.BlockRef, error) {
	d := DirHandle{
		BlockStore: fs.BlockStore}
	if err := d.Load(ref); err != nil {
		return nil, err
	}
	return d.Entries(), nil
}

func (fs FileStore) Grabs(ref base.BlockRef) base.BlockRefStream {
	reader := GrabReader{
		BlockStore: fs.BlockStore,
		ch:         make(base.BlockRefStream, 10),
	}

	go func() {
		reader.traceRef(ref)
		close(reader.ch)
	}()

	return reader.ch
}
