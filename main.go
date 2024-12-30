package main

import (
	"fmt"
	"log"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

const (
	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"
	//	filename string = "go-nebulon"
	tempfile string = "test1.tmp"
)

var fs base.FileStore

func getFile(fn string, ref wire.BlockRef) {
	reader, headers, err := fs.ReadStream(ref)

	if headers != nil {
		log.Printf("got headers\n")
		log.Printf("--> %s\n", headers)
	} else {
		log.Printf("NO HEADERS\n")
	}

	if err != nil {
		panic(fmt.Sprintf("open reader failed: %s", err))
	}

	err = util.CopyStreamToFile(reader, fn)
	if err != nil {
		panic(fmt.Errorf("copy failed [%w]", err))
	}
}

func main() {
	fs = filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

	log.Printf("Storing file: %s\n", filename)
	ref, err := helpers.StoreFile(fs, map[string]string{
		"wurst": "brot"}, filename)

	if err != nil {
		panic(fmt.Sprintf("ERROR %s\n", err))
	}

	log.Printf("Stored file ref %s\n", ref.Dump())
	getFile(tempfile, ref)
	log.Printf("Pulled file\n")
}
