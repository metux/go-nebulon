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

func (fs FileStore) writeBlockRefList(reflist wire.BlockRefList) (wire.BlockRef, error) {
	data, err := reflist.Marshal()

	if err != nil {
		log.Println("marshal error: ", err)
		return wire.BlockRef{}, err
	}

	key, encrypted, err := blockcrypt.BlockEncrypt(fs.encryption, data)
	if err != nil {
		log.Printf("storeDataBlock: BlockEncrypt() error %s\n", err)
		return wire.BlockRef{}, err
	}

	ref, err := fs.BlockStore.StoreBlock(encrypted)
	if err != nil {
		log.Println("error storing reflist block", err)
		return ref, err
	}

	ref.Type = wire.RefType_RefList
	ref.Cipher = fs.encryption
	ref.Key = key
	return ref, nil
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

func (fs FileStore) storeDataBlock(data []byte) (wire.BlockRef, error) {
	key, encrypted, err := blockcrypt.BlockEncrypt(fs.encryption, data)
	if err != nil {
		log.Printf("storeDataBlock: BlockEncrypt() error %s\n", err)
		return wire.BlockRef{}, err
	}

	ref, err := fs.BlockStore.StoreBlock(encrypted)
	if err != nil {
		log.Printf("storeDataBlock: storing block failed %s\n", err)
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

func (fs FileStore) storeFileData(r io.Reader) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}
	buf := make([]byte, BlockSize)
	for {
		readTotal, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				return reflist, err
			}
			break
		}
		ref, err := fs.storeDataBlock(buf[:readTotal])
		if err != nil {
			return reflist, err
		}
		reflist.AddRef(ref)
	}
	return reflist, nil
}

func (fs FileStore) storeRefLists(reflist wire.BlockRefList) (wire.BlockRef, error) {
	if reflist.Count() <= BlockListMax {
		return fs.writeBlockRefList(reflist)
	}

	new_reflist := wire.BlockRefList{}
	for reflist.Count() > BlockListMax {
		sub := wire.BlockRefList{Refs: reflist.Refs[:BlockListMax]}
		subref, err := fs.writeBlockRefList(sub)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
		reflist.Refs = reflist.Refs[BlockListMax:]
	}

	if reflist.Count() > 0 {
		subref, err := fs.writeBlockRefList(reflist)
		if err != nil {
			return wire.BlockRef{}, err
		}
		new_reflist.AddRef(subref)
	}

	return fs.storeRefLists(new_reflist)
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
	reflist, err := fs.storeFileData(r)
	if err != nil {
		return wire.BlockRef{}, err
	}

	content_ref, err := fs.storeRefLists(reflist)
	if err != nil {
		log.Printf("storeRefLists() error %s\n", err)
		return wire.BlockRef{}, err
	}

	log.Printf("StoreFile: Content ref: %s\n", content_ref.Dump())

	// fixme: add headers
	key, encrypted, err := fs.encodeFileControl(content_ref)

	if err != nil {
		return content_ref, nil
	}

	// fixme: must drop the key from content the ref !
	filehead := wire.FileHead{
// need to create public graph
//		Content: &content_ref,
		Private: encrypted,
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
		Oid: ref.Oid,
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
