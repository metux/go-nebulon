package blockstore

import (
	"log"
	"os"
	"path/filepath"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

type Store struct {
	Path string
}

func NewStore(path string) base.BlockStore {
	return Store{
		Path: path,
	}
}

func (s Store) Ref2FN(ref wire.BlockRef) string {
	return s.Path + "/" + ref.HexKey()
}

func (s Store) StoreBlock(data []byte) (wire.BlockRef, error) {
	ref := wire.RefForBlock(data)
	fn := s.Ref2FN(ref)
	log.Printf("Storing block as: %s\n", fn)
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)
	err := os.WriteFile(fn, data, 0644)
	return ref, err
}

func (s Store) LoadBlock(ref wire.BlockRef) ([]byte, error) {
	return os.ReadFile(s.Ref2FN(ref))
}
