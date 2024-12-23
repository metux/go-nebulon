package blockstore

import (
	"os"
	"path/filepath"
	"fmt"
)

type Store struct {
	Path string
}

func NewStore(path string) Store {
	return Store{
		Path: path,
	}
}

func (s Store) OID2FN(k OID) string {
	return s.Path + "/" + k.String()
}

func (s Store) StoreRaw(data []byte) OID {
	k := OIDForData(data)
	fn := s.OID2FN(k)
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)
	os.WriteFile(fn, data, 0644)
	fmt.Println("file name", fn)
	return k
}

func (s Store) LoadRaw(k OID) ([] byte, error) {
	return os.ReadFile(s.OID2FN(k))
}
