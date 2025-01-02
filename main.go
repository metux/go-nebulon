package main

import (
	"log"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/httpd"
	"github.com/metux/go-nebulon/util"
)

const (
	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"

// filename string = "go-nebulon"
)

func runServer(fs base.FileStore) {
	srv := httpd.NewServer(fs)
	srv.DoUpload(filename, "video/mp4")
	log.Printf("UP: %s\n", srv.Ref.Dump())
	srv.Run(":8080")
}

func main() {
	fs := filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

	ref, err := helpers.PutDirectory(fs, "", ".", util.FilterSkipHidden)
	if err != nil {
		panic(err)
	}

	log.Printf("Dir ref %s\n", ref.Dump())

	dir, err := fs.OpenDirectory(ref)

	if err != nil {
		panic(err)
	}

	log.Printf("DIR %+v\n", dir)

	// runServer(fs)
}
