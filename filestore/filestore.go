package filestore

import (
	"fmt"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
	"io"
	"log"
)

const (
	BlockSize    = 4096 * 16
	BlockListMax = 4
)

type FileStore struct {
	BlockStore base.BlockStore
}

func NewFileStore(bs base.BlockStore) base.FileStore {
	return FileStore{
		BlockStore: bs,
	}
}

func (fs FileStore) writeBlockRefList(reflist wire.BlockRefList) (wire.BlockRef, error) {
	data, err := reflist.Marshal()

	if err != nil {
		log.Println("marshal error: ", err)
		return wire.BlockRef{}, err
	}

	oid, err := fs.BlockStore.StoreBlock(data)
	if err != nil {
		log.Println("error storing reflist block", err)
		return oid, err
	}

	oid.Type = wire.RefType_RefList
	return oid, err
}

// FIXME: support encryption
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

func (fs FileStore) StoreFileData(r io.Reader) (wire.BlockRefList, error) {
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
		ref, err := fs.StoreBlock(buf[:readTotal])
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

func (fs FileStore) StoreFile(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	reflist, err := fs.StoreFileData(r)
	if err != nil {
		return wire.BlockRef{}, err
	}

	ref, err := fs.storeRefLists(reflist)
	if err != nil {
		log.Printf("storeRefLists() error %s\n", err)
		return wire.BlockRef{}, err
	}

	return ref, nil
}

func (fs FileStore) ReadFile(ref wire.BlockRef) (io.Reader, map[string]string, error) {
	// load the index block
	_, err := fs.LoadBlock(ref)
	if err != nil {
		return nil, nil, err
	}

	reader := &FileReader{
		fs: fs,
	}

	if err = reader.AddRef(ref); err != nil {
		return nil, nil, err
	}

	return reader, nil, nil
}
