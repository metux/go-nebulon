package filestore

import (
	"fmt"
	"io"

	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/blockcrypt"
	"github.com/metux/go-nebulon/wire"
)

type fileReader struct {
	util.ChainedReader
	fs FileStore
}

func (reader *fileReader) AddRef(ref wire.BlockRef) error {
	data, err := reader.LoadBlock(ref)
	if err != nil {
		return err
	}

	switch ref.Type {
	case wire.RefType_Blob:
		reader.AddBytes(data)
	case wire.RefType_RefList:
		bl, err := reader.loadBlockList(ref)
		if err != nil {
			return err
		}
		for _, walk := range bl.Refs {
			reader.AddRef(*walk)
		}
	default:
		return fmt.Errorf("unsupported ref type %+v\n", ref.Type)
	}

	return nil
}

func (r * fileReader) ReadStream(ref wire.BlockRef) (io.Reader, map[string]string, error) {
	// load the index block -- strip off the, that's later used used to decrypt the FileControl block
	filehead_ref := wire.BlockRef{
		Oid:  ref.Oid,
		Type: ref.Type,
	}

	filehead_bin, err := r.LoadBlock(filehead_ref)
	if err != nil {
		return nil, nil, fmt.Errorf("failed loading FileHead [%w]", err)
	}
	filehead := wire.FileHead{}
	filehead.Unmarshal(filehead_bin)

	fctrl_bin, err := blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, filehead.Private)
	if err != nil {
		return nil, nil, fmt.Errorf("error decrypting FileControl [%w]", err)
	}

	fctrl := wire.FileControl{}
	fctrl.Unmarshal(fctrl_bin)

	if err = r.AddRef(*fctrl.Content); err != nil {
		return nil, nil, err
	}

	return r, fctrl.Headers, nil
}

func (r * fileReader) loadBlockList(ref wire.BlockRef) (wire.BlockRefList, error) {
	reflist := wire.BlockRefList{}

	encrypted, err := r.fs.BlockStore.LoadBlock(ref)
	if err != nil {
		return reflist, fmt.Errorf("failed loading blocklist block [%w]", err)
	}

	data, err := blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, encrypted)
	if err != nil {
		return reflist, fmt.Errorf("failed decrypting blocklist [%w]", err)
	}

	// note do it in separate steps, since reflist is changed here
	err = reflist.Unmarshal(data)
	return reflist, err
}

func (r * fileReader) LoadBlock(ref wire.BlockRef) ([]byte, error) {
	data, err := r.fs.BlockStore.LoadBlock(ref)
	if err != nil {
		return data, err
	}

	return blockcrypt.BlockDecrypt(ref.Cipher, ref.Key, data)
}
