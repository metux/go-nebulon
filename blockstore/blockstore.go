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

func (s Store) Key2FN(k Key) string {
	return s.Path + "/" + k.String()
}

func (s Store) StoreRaw(data []byte) Key {
	k := KeyForData(data)
	fn := s.Key2FN(k)
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)
	os.WriteFile(fn, data, 0644)
	fmt.Println("file name", fn)
	return k
}

func (s Store) LoadRaw(k Key) ([] byte, error) {
	fn := s.Key2FN(k)
	return os.ReadFile(fn)
}
