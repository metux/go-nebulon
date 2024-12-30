package filestore

import (
	"io"
	"log"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

const (
	//	BlockSize    = 4096 * 16
	BlockSize    = 4096 * 1024
	BlockListMax = BlockSize / 80 // a blocklist entry is about 80 bytes
)

type FileStore struct {
	BlockStore base.BlockStore
	encryption wire.CipherType
}

func NewFileStore(bs base.BlockStore) base.FileStore {
	return FileStore{
		BlockStore: bs,
		encryption: wire.CipherType_AES_CBC_ZSTD,
	}
}

func (fs FileStore) loadBlockList(ref wire.BlockRef) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}

	encrypted, err := fs.BlockStore.LoadBlock(ref)
	if err != nil {
		log.Printf("failed loading blocklist block: %s\n", err)
		return reflist, err
	}

	data, err := blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, encrypted)
	if err != nil {
		log.Printf("failed decrypting blocklist: %s\n", err)
		return reflist, err
	}

	// note do it in separate steps, since reflist is changed here
	err = reflist.Unmarshal(data)
	return reflist, err
}

func (fs FileStore) writeDataBlock(data []byte) (wire.BlockRef, error) {
	key, encrypted, err := blockcrypt.BlockEncrypt(fs.encryption, data)
	if err != nil {
		log.Printf("writeDataBlock: BlockEncrypt() error %s\n", err)
		return wire.BlockRef{}, err
	}

	ref, err := fs.BlockStore.StoreBlock(encrypted)
	if err != nil {
		log.Printf("writeDataBlock: storing block failed %s\n", err)
		return wire.BlockRef{}, err
	}

	ref.Key = key
	ref.Cipher = fs.encryption

	return ref, nil
}

func (fs FileStore) LoadBlock(ref wire.BlockRef) ([]byte, error) {
	log.Printf("LoadBlock oid=%s:%s:%X key=%X\n", ref.Type, ref.Cipher, ref.Oid, ref.Key)

	data, err := fs.BlockStore.LoadBlock(ref)
	if err != nil {
		return data, err
	}

	return blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, data)
}

func (fs FileStore) encodeFileControl(content_ref wire.BlockRef, headers map[string]string) ([]byte, []byte, error) {

	// fixme: add headers
	fctrl := wire.FileControl{
		Content: &content_ref,
		Headers: headers,
	}
	data, err := fctrl.Marshal()
	if err != nil {
		log.Printf("failed marshalling fctrl: %s\n", err)
		return []byte{}, []byte{}, err
	}

	key, encrypted, err := blockcrypt.BlockEncrypt(fs.encryption, data)
	if err != nil {
		log.Printf("failed encrypting fctrl: %s\n", err)
		return []byte{}, []byte{}, err
	}

	return key, encrypted, nil
}

func (fs FileStore) StoreStream(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	context := fileWriteContext{
		fs: fs,
	}

	content_ref, err := context.storeFileStream(r)

	log.Printf("StoreFile: Content ref: %s\n", content_ref.Dump())

	// fixme: add headers
	key, encrypted, err := fs.encodeFileControl(content_ref, headers)

	if err != nil {
		return wire.BlockRef{}, err
	}

	graph_ref, err := context.writeGraph()
	if err != nil {
		return wire.BlockRef{}, err
	}

	filehead_ref, err := context.writeFileHead(encrypted, graph_ref)
	if err != nil {
		log.Printf("error storing file head in blockstore %s\n", err)
		return content_ref, err
	}

	log.Printf("file head ref: %X\n", filehead_ref.Oid)

	filehead_ref.Cipher = fs.encryption
	filehead_ref.Key = key
	filehead_ref.Type = wire.RefType_File
	return filehead_ref, nil
}

func (fs FileStore) ReadStream(ref wire.BlockRef) (io.Reader, map[string]string, error) {
	// load the index block
	filehead_ref := wire.BlockRef{
		Oid:  ref.Oid,
		Type: ref.Type,
	}

	filehead_bin, err := fs.LoadBlock(filehead_ref)
	if err != nil {
		log.Printf("failed loading filehead block: %s\n", err)
		return nil, nil, err
	}
	filehead := wire.FileHead{}
	filehead.Unmarshal(filehead_bin)

	fctrl_bin, err := blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, filehead.Private)
	if err != nil {
		log.Printf("error decrypting fctrl: %s\n", err)
		return nil, nil, err
	}

	fctrl := wire.FileControl{}
	fctrl.Unmarshal(fctrl_bin)

	log.Printf("ReadFile content %s\n", fctrl.Content.Dump())

	reader := &FileReader{
		fs: fs,
	}

	if err = reader.AddRef(*fctrl.Content); err != nil {
		return nil, nil, err
	}

	return reader, nil, nil
}
