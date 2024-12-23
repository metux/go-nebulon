package filestore

import (
    "errors"
    "bytes"
    "io"
    "fmt"
    "github.com/metux/go-nebulon/base"
    "github.com/metux/go-nebulon/wire"
//    "google.golang.org/protobuf/proto"
)

type FileStore struct {
	BlockStore base.BlockStore
}

func NewFileStore(bs base.BlockStore) base.FileStore {
    return FileStore {
	BlockStore: bs,
    }
}

func (fs FileStore) StoreBlockList(oids [] base.OID) (base.OID, error) {
//	refs := wire.EncapOIDRefList(oids)

	fmt.Println("OIDS to store", oids)
	data, err := wire.MarshalOIDRefList(oids)

	if err != nil {
		fmt.Println("marshal error: ", err)
	}

	fmt.Println(data)

	return base.OID{}, nil
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
			if !bytes.Equal(d, buf[:readTotal]) {
				fmt.Printf("Read back failed - blocks not equal\n")
				return base.OID{}, errors.New("Read back failed - blocks not equal")
			}
		}

		oids = append(oids, k)
	}

	return fs.StoreBlockList(oids)
}

func (fs FileStore) ReadFile(oid base.OID) (io.Reader, map[string]string, error) {
    return nil, nil, nil
}
