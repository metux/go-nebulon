package filestore

import (
	"fmt"
	"io"
	"log"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

const (
	//	BlockSize    = 4096 * 16
	BlockSize     = 4096 * 1024
	BlockListMax  = BlockSize / 80 // a blocklist entry is about 80 bytes
	DefaultCipher = wire.CipherType_AES_CBC_ZSTD
)

type FileStore struct {
	BlockStore base.BlockStore
	encryption wire.CipherType
}

func NewFileStore(bs base.BlockStore) base.FileStore {
	return FileStore{
		BlockStore: bs,
		encryption: DefaultCipher,
	}
}

func (fs FileStore) loadBlockList(ref wire.BlockRef) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}

	encrypted, err := fs.BlockStore.LoadBlock(ref)
	if err != nil {
		return reflist, fmt.Errorf("failed loading blocklist block [%w]", err)
	}

	data, err := blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, encrypted)
	if err != nil {
		return reflist, fmt.Errorf("failed decrypting blocklist [%w]", err)
	}

	// note do it in separate steps, since reflist is changed here
	err = reflist.Unmarshal(data)
	return reflist, err
}

func (fs FileStore) LoadBlock(ref wire.BlockRef) ([]byte, error) {
	log.Printf("LoadBlock oid=%s:%s:%X key=%X\n", ref.Type, ref.Cipher, ref.Oid, ref.Key)

	data, err := fs.BlockStore.LoadBlock(ref)
	if err != nil {
		return data, err
	}

	return blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, data)
}

func (fs FileStore) StoreStream(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	context := fileWriteContext{
		fs:           fs,
		cipher:       fs.encryption,
		blockSize:    BlockSize,
		blockListMax: BlockListMax,
	}
	return context.StoreStream(r, headers)
}

func (fs FileStore) ReadStream(ref wire.BlockRef) (io.Reader, map[string]string, error) {
	reader := &fileReader{
		fs: fs,
	}

	return reader.ReadStream(ref)
}
