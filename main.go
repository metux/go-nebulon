package main

import (
	"fmt"
	"log"
	"io/ioutil"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/httpd"
)

const (
	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"
//	filename string = "go-nebulon"
)

var fs base.FileStore

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

func main() {
	fs = filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

//	storeDirectory(fs, ".")

	srv := httpd.NewServer(fs)
	srv.DoUpload(filename, "video/mp4")
	log.Printf("UP: %s\n", srv.Ref.Dump())
	srv.Run(":8080")
}
