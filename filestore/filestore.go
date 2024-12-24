package filestore

import (
	"fmt"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
	"google.golang.org/protobuf/proto"
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
	fmt.Println("numer of OIDs", len(refs))

	reflist := wire.BlockRefList{
		Magic: "BLOCK REF LIST",
		Refs:  refs,
	}

	data, err := proto.Marshal(&reflist)

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

func (fs FileStore) LoadBlockList(ref wire.BlockRef) ([]*wire.BlockRef, error) {
	return make([]*wire.BlockRef, 0), nil
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

	return nil, nil, nil
}
