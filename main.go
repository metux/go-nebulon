package main

import (
	"fmt"
	"github.com/metux/go-nebulon/blockstore"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/helpers"
)

func main() {
	fs := filestore.NewFileStore(blockstore.NewStore(".storedata"))
	oid, err := helpers.StoreFile(fs, map[string]string{}, "go-nebulon")

	if err != nil {
		fmt.Printf("ERROR %s\n", err)
	}
	fmt.Printf("OID %s\n", oid)
}
