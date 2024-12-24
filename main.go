package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
)

var fs base.FileStore

func getFile(fn string, ref wire.BlockRef) {

	reader, headers, err := fs.ReadFile(ref)
	log.Printf("opened reader\n")

	if headers != nil {
		log.Printf("got headers\n")
	}
	if err != nil {
		panic(fmt.Sprintf("open reader failed: %s", err))
	}

	log.Printf("got reader\n")
	newf, err := os.Create(fn)
	if err != nil {
		panic(fmt.Sprintf("open write temp file failed: %s", err))
	}
	defer newf.Close()

	buf := make([]byte, 128)
	for {
		readTotal, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		log.Printf("writing %d bytes\n", readTotal)

		n, err := newf.Write(buf[:readTotal])
		if err != nil {
			log.Printf("writing failed: %s\n", err)
		}
		log.Printf("written: %d\n", n)
	}
	log.Printf("finished\n");
}

func main() {
	fs = filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))
	ref, err := helpers.StoreFile(fs, map[string]string{}, "go-nebulon")

	if err != nil {
		panic(fmt.Sprintf("ERROR %s\n", err))
	}

	log.Printf("Stored file ref %s\n", ref.HexKey())
	getFile("go-nebulon.tmp", ref)
}
