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

func (fs FileStore) encodeFileControl(content_ref wire.BlockRef) ([]byte, []byte, error) {

	// fixme: add headers
	fctrl := wire.FileControl{
		Content: &content_ref,
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

func (fs FileStore) StoreFile(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	context := fileWriteContext{
		fs: fs,
	}

	reflist, err := context.storeFileData(r)
	if err != nil {
		return wire.BlockRef{}, err
	}

	content_ref, err := context.storeRefLists(reflist, context.fs.encryption)
	if err != nil {
		log.Printf("storeRefLists() error %s\n", err)
		return wire.BlockRef{}, err
	}

	log.Printf("StoreFile: Content ref: %s\n", content_ref.Dump())

	// fixme: add headers
	// fixme: store context graph
	key, encrypted, err := fs.encodeFileControl(content_ref)

	if err != nil {
		return wire.BlockRef{}, err
	}

	context.graph.Sort()
	graph_ref, err := context.storeRefLists(context.graph, wire.CipherType_None)
	log.Printf("Graph ref: %s\n", graph_ref.Dump())

	// fixme: must drop the key from content the ref !
	filehead := wire.FileHead{
		Private: encrypted,
		Graph: &graph_ref,
	}
	filehead_bin, err := filehead.Marshal()
	if err != nil {
		log.Printf("error marshalling file head: %s\n", err)
		return content_ref, err
	}

	filehead_ref, err := fs.BlockStore.StoreBlock(filehead_bin)
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

func (fs FileStore) ReadFile(ref wire.BlockRef) (io.Reader, map[string]string, error) {
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
