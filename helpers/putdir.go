package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/util"
)

var (
	TracePutDirectory = false
)

func PutDirectory(fs base.IFileStore, name string, dirname string, filter util.FileNameFilter) (base.BlockRef, error) {
	if TracePutDirectory {
		util.TimeTrack(time.Now(), "StoreDirectory ("+dirname+")")
	}
	items, _ := ioutil.ReadDir(dirname)
	refEntries := base.BlockRefList{}
	for _, item := range items {
		name := item.Name()
		if !filter(name, dirname) {
			continue
		}

		fn := filepath.Clean(dirname + string(os.PathSeparator) + name)
		if item.IsDir() {
			ref, err := PutDirectory(fs, name, fn, filter)
			if err != nil {
				log.Printf("error storing subdir %s\n", err)
				return ref, err
			}
			if TracePutDirectory {
				log.Printf("Stored directory %s\n", ref.Dump())
			}
			refEntries.Add(ref)
		} else {
			// handle file there
			ref, err := PutFile(fs, name, wire.Header{}, fn)
			if err != nil {
				return ref, fmt.Errorf("error storing file [%w]\n", err)
			}
			if TracePutDirectory {
				log.Printf("Stored file %s\n", ref.Dump())
			}
			refEntries.Add(ref)
		}
	}

	ref, err := fs.StoreDirectory(refEntries)
	ref.Name = name
	return ref, err
}
