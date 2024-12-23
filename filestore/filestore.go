package filestore

import (
    "errors"
    "bytes"
    "io"
    "fmt"
    "github.com/metux/go-nebulon/base"
)

type FileStore struct {
	BlockStore base.BlockStore
}

func NewFileStore(bs base.BlockStore) base.FileStore {
    return FileStore {
	BlockStore: bs,
    }
}

func (fs FileStore) StoreFile(r io.Reader, headers map[string]string) (base.OID, error) {
	// declare chunk size
	const blocksize = 4096

	oids := make([]base.OID, 1)

	buf := make([]byte, blocksize)
	for {
		readTotal, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				return base.OID{}, err
			}
			break
		}
//              fmt.Println(string(b[:readTotal])) // print content from buffer
		k,_ := fs.BlockStore.StoreBlock(buf[:readTotal])
		for _, v := range k.Data {
			fmt.Printf("%d ", v)
		}
		fmt.Println("\n")

		d,e := fs.BlockStore.LoadBlock(k)
		if e != nil {
			fmt.Printf("Read back error %s\n", e)
			return base.OID{}, e
		} else {
			if bytes.Equal(d, buf[:readTotal]) {
				fmt.Printf("Read back OK\n")
			} else {
				fmt.Printf("Read back failed - blocks not equal\n")
				return base.OID{}, errors.New("Read back failed - blocks not equal")
			}
		}

		oids = append(oids, k)
	}

	for _,o := range oids {
		fmt.Printf("OID: %s\n", o.String())
	}

	return base.OID{}, nil
}

func (fs FileStore) ReadFile(oid base.OID) (io.Reader, map[string]string, error) {
    return nil, nil, nil
}
