package filestore

import (
    "errors"
    "bytes"
    "io"
    "fmt"
    "github.com/metux/go-nebulon/base"
    "github.com/metux/go-nebulon/wire"
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

func (fs FileStore) StoreBlockList(oids [] base.BlockRef) (base.OID, error) {
	// FIXME: should split large chunks

	fmt.Println("OIDS to store", oids)
	fmt.Println("numer of OIDs", len(oids))
	data, err := wire.MarshalOIDRefList(oids)

	if err != nil {
		fmt.Println("marshal error: ", err)
		return base.OID{}, err
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

func (fs FileStore) StoreFile(r io.Reader, headers map[string]string) (base.OID, error) {
	oids := make([]base.BlockRef, 1)

	buf := make([]byte, BlockSize)
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

		oids = append(oids, base.BlockRef{OID: k, Type: base.RefType_Blob})
	}

	return fs.StoreBlockList(oids)
}

func (fs FileStore) ReadFile(oid base.OID) (io.Reader, map[string]string, error) {
    return nil, nil, nil
}
