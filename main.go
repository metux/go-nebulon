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

func storeDirectory(fs base.FileStore, dir string, prefix string) (wire.BlockRef, error) {
	log.Printf("Directory: %s %s\n", prefix, dir)
	items, _ := ioutil.ReadDir(dir)
	refEntries := wire.BlockRefList{}
	for _, item := range items {
		name := item.Name()
		if name[0] == '.' {
			continue
		}

		fn := appendDir(dir, name)
		if item.IsDir() {
			ref, err := storeDirectory(fs, fn, prefix+"/"+dir)
			if err != nil {
				return ref, err
			}
			log.Printf("stored directory %s\n", ref.Dump())
		} else {
			// handle file there
			ref, err := helpers.StoreFile(fs, name, wire.Header{}, fn)
			if err != nil {
				return wire.BlockRef{}, fmt.Errorf("error storing file [%w]\n", err)
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

	storeDirectory(fs, ".", "")
	// runServer(fs)
}
