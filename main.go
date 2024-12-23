package main

import (
	"fmt"
	//    "path"
	"github.com/metux/go-nebulon/blockstore"
	"io"
	"os"
)

func storefile(st blockstore.Store, fn string) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()

	// declare chunk size
	const maxSz = 4096

	// create buffer
	b := make([]byte, maxSz)

	for {
		// read content to buffer
		readTotal, err := file.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Println(string(b[:readTotal])) // print content from buffer
		k := st.StoreRaw(b[:readTotal])
		for _, v := range k.Data {
			fmt.Printf("%d ", v)
		}
		fmt.Println("\n")
	}
	return nil
}

func main() {
	st := blockstore.NewStore(".storedata")
	storefile(st, "Makefile")
}
