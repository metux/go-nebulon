package main

import (
	"fmt"
	"log"
	"io/ioutil"

	"github.com/udhos/equalfile"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

const (
	//	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"
	filename string = "go-nebulon"
	tempfile string = "test1.tmp"
)

var fs base.FileStore

func getFile(fn string, ref wire.BlockRef) {
	reader, headers, err := fs.ReadStream(ref)
	log.Printf("Headers: %s\n", headers)

	if err != nil {
		panic(fmt.Sprintf("open reader failed [%s]", err))
	}

	err = util.CopyStreamToFile(reader, fn)
	if err != nil {
		panic(fmt.Errorf("copy failed [%w]", err))
	}
}

func appendDir(dn string, fn string) string {
	if dn == "." || dn == "" {
		return fn
	}
	return dn+"/"+fn
}

func storeDirectory(fs base.FileStore, dir string) {
	items, _ := ioutil.ReadDir(dir)
	for _, item := range items {
		name := item.Name()
		if name[0] == '.' {
			continue
		}

		fn := appendDir(dir, name)
		if item.IsDir() {
			fmt.Println("DIR: "+fn)
			storeDirectory(fs, fn)
		} else {
			// handle file there
			fmt.Println("FIL: "+fn)

			ref, err := helpers.StoreFile(fs, map[string]string{"filename":fn}, filename)
			if err != nil {
				log.Printf("error storing file [%s]\n", err)
			} else {
				log.Printf("stored %s\n", ref.Dump())
			}
		}
	}
}

func testFile(fs base.FileStore) {
	log.Printf("Storing file: %s\n", filename)
	ref, err := helpers.StoreFile(fs, map[string]string{
		"Content-Type": "video/mp4"}, filename)

	if err != nil {
		panic(err)
	}

	log.Printf("Stored file ref %s\n", ref.Dump())
	headers, err := helpers.GetFile(fs, tempfile, ref)
	if err != nil {
		panic(err)
	}
	log.Printf("Pulled file: headers=%s\n", headers)

	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
	equal, err := cmp.CompareFile(filename, tempfile)

	if equal {
		log.Printf("Both files are equal [%s]\n", err)
	} else {
		log.Printf("Files mismatch [%s]\n", err)
	}
}

func main() {
	fs = filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

	testFile(fs)
//	storeDirectory(fs, ".")
}
