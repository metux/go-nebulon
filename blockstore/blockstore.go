package blockstore

import (
	"os"
	"path/filepath"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
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
	err := os.WriteFile(fn, data, 0644)
	return ref, err
}

func (s SimpleStore) LoadBlock(ref wire.BlockRef) ([]byte, error) {
	return os.ReadFile(s.Ref2FN(ref))
}
