package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/httpd"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

const (
	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"

// filename string = "go-nebulon"
)

var fs base.FileStore

func appendDir(dn string, fn string) string {
	if dn == "." || dn == "" {
		return fn
	}
	return dn + "/" + fn
}

func fnFilter(name string, path string) bool {
	return true
}

func storeDirectory(fs base.FileStore, dirname string, filter util.FileNameFilter) (wire.BlockRef, error) {
	log.Printf("Directory: %s\n", dirname)
	items, _ := ioutil.ReadDir(dirname)
	refEntries := wire.BlockRefList{}
	for _, item := range items {
		name := item.Name()
		log.Printf("File entry: %s\n", name)
		if util.PathIsSelf(name) {
			log.Printf("skipping self")
			continue
		}
		if !filter(name, dirname) {
			log.Printf("skipped file")
			continue
		}

		fn := appendDir(dirname, name)
		if item.IsDir() {
			ref, err := storeDirectory(fs, fn, filter)
			if err != nil {
				return ref, err
			}
			log.Printf("stored directory %s\n", ref.Dump())
		} else {
			// handle file there
			ref, err := helpers.StoreFile(fs, name, wire.Header{}, fn)
			if err != nil {
				return ref, fmt.Errorf("error storing file [%w]\n", err)
			}
			refEntries.Add(ref)
		}
	}

	log.Printf("Entries: %d\n", len(refEntries.Refs))
	for idx, ent := range refEntries.Refs {
		log.Printf("[%d] %s\n", idx, ent.Dump())
	}

	return fs.StoreDirectory(refEntries)
}

func runServer(fs base.FileStore) {
	srv := httpd.NewServer(fs)
	srv.DoUpload(filename, "video/mp4")
	log.Printf("UP: %s\n", srv.Ref.Dump())
	srv.Run(":8080")
}

func main() {
	fs = filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

	storeDirectory(fs, ".", util.FilterSkipHidden)
	// runServer(fs)
}
