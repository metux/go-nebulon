package main

import (
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
	// filename string = "go-nebulon"
	filename string = "/home/nekrad/dl/000.capture/elen0_tg/elen_cross-2024-09-04-04-26-03.P00.mkv.00.mux.mp4.tg.mp4"

	test_temp_file = ".download.tmp"
	server_port    = ":8080"
)

func runServer(fs base.FileStore) {
	srv := httpd.NewServer(fs)
	srv.DoUpload(filename, wire.ContentType_MP4)
	log.Printf("UP: %s\n", srv.Ref.Dump())
	srv.Run(server_port)
}

func panicX(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fs := filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

	ref, err := helpers.PutDirectory(fs, "", ".", util.FilterSkipHidden)
	panicX(err)

	panicX(helpers.CompareTree(fs, ".", ref, test_temp_file))
	log.Printf("CompareTree() done\n")

	// runServer(fs)
}
