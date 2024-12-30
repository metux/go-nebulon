package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/wire"
)

const (
	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"
	//	filename string = "go-nebulon"
	tempfile string = "test1.tmp"
)

var fs base.FileStore

func getFile(fn string, ref wire.BlockRef) {
	reader, headers, err := fs.ReadFile(ref)

	if headers != nil {
		log.Printf("got headers\n")
	}
	if err != nil {
		panic(fmt.Sprintf("open reader failed: %s", err))
	}

	newf, err := os.Create(fn)
	if err != nil {
		panic(fmt.Sprintf("open write temp file failed: %s", err))
	}
	defer newf.Close()

	buf := make([]byte, 1024)
	for {
		readTotal, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		_, err = newf.Write(buf[:readTotal])
		if err != nil {
			log.Printf("writing failed: %s\n", err)
		}
	}
}

func main() {
	fs = filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

	log.Printf("Storing file: %s\n", filename)
	ref, err := helpers.StoreFile(fs, map[string]string{
		"wurst": "brot" }, filename)

	if err != nil {
		panic(fmt.Sprintf("ERROR %s\n", err))
	}

	log.Printf("Stored file ref %s\n", ref.Dump())
	getFile(tempfile, ref)
	log.Printf("Pulled file\n")
}
