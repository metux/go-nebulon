package blockstore

import (
	"os"
	"path/filepath"

	"github.com/metux/go-nebulon/base"
)

type Store struct {
	Path string
}

func NewStore(path string) base.BlockStore {
	return Store{
		Path: path,
	}
}

func (s Store) OID2FN(k base.OID) string {
	return s.Path + "/" + k.String()
}

func (s Store) StoreBlock(data []byte) (base.OID, error) {
	k := base.OIDForBlock(data)
	fn := s.OID2FN(k)
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)
	err := os.WriteFile(fn, data, 0644)
	return k, err
}

func (s Store) LoadBlock(k base.OID) ([] byte, error) {
	return os.ReadFile(s.OID2FN(k))
}
