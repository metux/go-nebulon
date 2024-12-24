package main

import (
	"log"
	"fmt"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
)

func main() {
	fs := filestore.NewFileStore(blockstore.NewSimpleStore(".storedata"))
	oid, err := helpers.StoreFile(fs, map[string]string{}, "go-nebulon")

	if err != nil {
		fmt.Printf("ERROR %s\n", err)
	}
	log.Printf("Stored file ref %s\n", oid.HexKey())

	reader, headers, err := fs.ReadFile(oid)

	if reader != nil {
		fmt.Printf("got reader\n")
	}
	if headers != nil {
		fmt.Printf("got headers\n")
	}
	if err != nil {
		fmt.Printf("ERR %s\n", err)
	}
}
