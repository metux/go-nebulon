package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

var (
	TracePutDirectory = false
)

func PutDirectory(fs base.FileStore, dirname string, filter util.FileNameFilter) (wire.BlockRef, error) {
	if TracePutDirectory {
		util.TimeTrack(time.Now(), "StoreDirectory ("+dirname+")")
	}
	items, _ := ioutil.ReadDir(dirname)
	refEntries := wire.BlockRefList{}
	for _, item := range items {
		name := item.Name()
		if !filter(name, dirname) {
			continue
		}

		fn := filepath.Clean(dirname + string(os.PathSeparator) + name)
		if item.IsDir() {
			ref, err := PutDirectory(fs, fn, filter)
			if err != nil {
				return ref, err
			}
			if TracePutDirectory {
				log.Printf("Stored directory %s\n", ref.Dump())
			}
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

	return fs.StoreDirectory(refEntries)
}
