package main

import (
	"bytes"
	"fmt"
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
//		fmt.Println(string(b[:readTotal])) // print content from buffer
		k := st.StoreRaw(b[:readTotal])
		for _, v := range k.Data {
			fmt.Printf("%d ", v)
		}
		fmt.Println("\n")

		d,e := st.LoadRaw(k)
		if e != nil {
			fmt.Printf("Read back error %s\n", e)
		} else {
			if bytes.Equal(d, b[:readTotal]) {
				fmt.Printf("Read back OK\n")
			} else {
				fmt.Printf("Read back failed - blocks not equal\n")
			}
		}
	}
	return nil
}

func main() {
	st := blockstore.NewStore(".storedata")
	storefile(st, "go-nebulon")
}
