package filestore

import (
	"fmt"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
	"io"
	"log"
)

const (
	BlockSize = 4096 * 16
)

type FileStore struct {
	BlockStore base.BlockStore
}

func NewFileStore(bs base.BlockStore) base.FileStore {
	return FileStore{
		BlockStore: bs,
	}
}

func (fs FileStore) StoreBlockList(refs []*wire.BlockRef) (wire.BlockRef, error) {
	// FIXME: should split large chunks

	fmt.Println("OIDS to store", refs)

	reflist := wire.BlockRefList{
		Refs: refs,
	}

	fmt.Println("number of ref entries", reflist.Count())
	data, err := reflist.Marshal()

	if err != nil {
		fmt.Println("marshal error: ", err)
		return wire.BlockRef{}, err
	}

	fmt.Println(data)

	oid, err := fs.BlockStore.StoreBlock(data)
	if err != nil {
		fmt.Println("error storing reflist block", err)
		return oid, err
	}

	return oid, err
}

func (fs FileStore) LoadBlockList(ref wire.BlockRef) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}

	data, err := fs.BlockStore.LoadBlock(ref)

	if err != nil {
		log.Printf("failed loading blocklist block: %s\n", err)
		return reflist, err
	}

	// note do it in separate steps, since reflist is changed here
	err = reflist.Unmarshal(data)
	return reflist, err
}

func (fs FileStore) StoreBlock(data []byte) (wire.BlockRef, error) {
	// FIXME: need to encrypt
	ref, err := fs.BlockStore.StoreBlock(data)
	return ref, err
}

func (fs FileStore) LoadBlock(ref wire.BlockRef) ([]byte, error) {
	// FIXME: need to decrypt
	data, err := fs.BlockStore.LoadBlock(ref)
	return data, err
}

func (fs FileStore) StoreFile(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	oids := make([]*wire.BlockRef, 1)

	buf := make([]byte, BlockSize)
	for {
		readTotal, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				return wire.BlockRef{}, err
			}
			break
		}
		ref, err := fs.StoreBlock(buf[:readTotal])
		if err != nil {
			log.Println("StoreBlock error", err)
			return wire.BlockRef{}, err
		}

		oids = append(oids, &ref)
	}

	return fs.StoreBlockList(oids)
}

func (fs FileStore) ReadFile(ref wire.BlockRef) (io.Reader, map[string]string, error) {
	// load the index block
	_, err := fs.LoadBlock(ref)
	if err != nil {
		log.Printf("ReadFile: error: %s\n", err)
		return nil, nil, err
	}

	bl, err := fs.LoadBlockList(ref)
	if err != nil {
		log.Printf("ReadFile: failed reading block list %s\n", err)
	}

	log.Printf("BLOCK REF LIST %+v\n", bl)
	return nil, nil, nil
}
