package blockstore

import (
	//	"io/ioutil"
	"os"
	"path/filepath"
	//    "crypto/sha256"
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
	fmt.Println("file name", fn)
	return k
}
