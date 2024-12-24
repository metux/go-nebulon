package filestore

import (
    "errors"
    "bytes"
    "io"
    "fmt"
    "github.com/metux/go-nebulon/base"
    "github.com/metux/go-nebulon/wire"
    "google.golang.org/protobuf/proto"
)

const (
	BlockSize = 4096 * 16
)

type FileStore struct {
	BlockStore base.BlockStore
}

func NewFileStore(bs base.BlockStore) base.FileStore {
    return FileStore {
	BlockStore: bs,
    }
}

func (fs FileStore) StoreBlockList(refs [] * wire.BlockRef) (wire.BlockRef, error) {
	// FIXME: should split large chunks

	fmt.Println("OIDS to store", refs)
	fmt.Println("numer of OIDs", len(refs))

	reflist := wire.BlockRefList{
		Magic: "BLOCK REF LIST",
		Refs: refs,
	}

	data, err := proto.Marshal(&reflist)

	if err != nil {
		fmt.Println("marshal error: ", err)
		return wire.BlockRef{}, err
	}

	fmt.Println(data)

	oid, err := fs.BlockStore.StoreBlock(data)
	if err != nil {
		fmt.Println("error storing reflist block", err)
		return oid, err
	}

	return oid, err
}

//func (fs FileStore) StoreBlock(date [] byte) (base.Ref) {
//    k,_ := fs.BlockStore.StoreBlock(data)
//}

func (fs FileStore) StoreFile(r io.Reader, headers map[string]string) (wire.BlockRef, error) {
	oids := make([]*wire.BlockRef, 1)

	buf := make([]byte, BlockSize)
	for {
		readTotal, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				return wire.BlockRef{}, err
			}
			break
		}
		ref,_ := fs.BlockStore.StoreBlock(buf[:readTotal])
		for _, v := range ref.Data {
			fmt.Printf("%d ", v)
		}
		fmt.Println("\n")

		d,e := fs.BlockStore.LoadBlock(ref)
		if e != nil {
			fmt.Printf("Read back error %s\n", e)
			return wire.BlockRef{}, e
		} else {
			if !bytes.Equal(d, buf[:readTotal]) {
				fmt.Printf("Read back failed - blocks not equal\n")
				return wire.BlockRef{}, errors.New("Read back failed - blocks not equal")
			}
		}

		oids = append(oids, &ref)
	}

	return fs.StoreBlockList(oids)
}

func (fs FileStore) ReadFile(oid base.OID) (io.Reader, map[string]string, error) {
    return nil, nil, nil
}
