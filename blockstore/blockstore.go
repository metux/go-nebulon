package blockstore

import (
	"log"
	"os"
	"path/filepath"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

var (
	TraceWrite = false
)

type SimpleStore struct {
	Path string
}

func NewSimpleStore(path string) base.BlockStore {
	return SimpleStore{
		Path: path,
	}
}

func (s SimpleStore) Ref2FN(ref wire.BlockRef) string {
	return s.Path + "/" + ref.OID()
}

func (s SimpleStore) StoreBlock(data []byte) (wire.BlockRef, error) {
	ref := wire.RefForBlock(data)
	fn := s.Ref2FN(ref)
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)

	// dont write already existing objects
	// 2do: we could check equality here
	if _, err := os.Stat(fn); err == nil {
		if TraceWrite {
			log.Printf("object already exists %s\n", fn)
		}
		return ref, nil
	}
	err := os.WriteFile(fn, data, 0644)
	return ref, err
}

func (s SimpleStore) LoadBlock(ref wire.BlockRef) ([]byte, error) {
	return os.ReadFile(s.Ref2FN(ref))
}
