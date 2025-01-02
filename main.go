package main

import (
	"fmt"
	"log"

	"github.com/udhos/equalfile"

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
	server_port = ":8080"
)

func runServer(fs base.FileStore) {
	srv := httpd.NewServer(fs)
	srv.DoUpload(filename, wire.ContentType_MP4)
	log.Printf("UP: %s\n", srv.Ref.Dump())
	srv.Run(server_port)
}


func compareEntry(fs base.FileStore, path string, entry wire.BlockRef) error {
	fn := path + "/" + entry.Name
	log.Printf("FILE: %s --- %s\n", fn, entry.Dump())

	if entry.IsFile() {
		_, err := helpers.GetFile(fs, test_temp_file, entry)
		if err != nil {
			return err
		}

		cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
		equal, err := cmp.CompareFile(test_temp_file, fn)

		if err != nil {
			return fmt.Errorf("compare error [%w]", err)
		}

		if !equal {
			return fmt.Errorf("files mismatch: %s\n", fn)
		}

		log.Printf("file %s OK\n", fn)
	} else if entry.IsDir() {
		log.Printf("skipping dir %s\n", fn)
	} else {
		return fmt.Errorf("unexpected object, neither dir nor file")
	}

	return nil
}

func compareEntries(fs base.FileStore, path string, entries []wire.BlockRef) error {
	log.Printf("PATH: %s\n", path)
	for _, e := range entries {
		if err := compareEntry(fs, path, e); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fs := filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))

	ref, err := helpers.PutDirectory(fs, "", ".", util.FilterSkipHidden)
	if err != nil {
		panic(err)
	}

	//	log.Printf("Dir ref %s\n", ref.Dump())

	entries, err := fs.ReadDirectory(ref)
	if err != nil {
		panic(err)
	}

	if err = compareEntries(fs, ".", entries); err != nil {
		panic(err)
	}

	// runServer(fs)
}
