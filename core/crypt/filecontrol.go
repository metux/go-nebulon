package crypt

import (
	"fmt"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

func EncryptFileControl(cipher base.CipherType, fctrl wire.FileControl) ([]byte, []byte, wire.CipherType, error) {
	data, err := fctrl.Marshal()
	if err != nil {
		return nil, nil, cipher, fmt.Errorf("failed marshalling FileControl [%w]", err)
	}

	key, encrypted, cipher, err := BlockEncrypt(cipher, data)
	if err != nil {
		return nil, nil, cipher, fmt.Errorf("failed encrypting FileControl [%w]", err)
	}

	return key, encrypted, cipher, nil
}

func DecryptFileControl(cipher base.CipherType, key []byte, private []byte) (wire.FileControl, error) {
	fctrl_bin, err := BlockDecrypt(cipher, key, private)
	if err != nil {
		return wire.FileControl{}, fmt.Errorf("error decrypting FileControl [%w]", err)
	}

	return wire.FileControlUnmarshal(fctrl_bin)
}

func LoadFileHead(bs base.IBlockStore, ref base.BlockRef) (wire.FileHead, error) {
	// load the index block -- strip off the key, that's later used used to decrypt the FileControl block
	// FIXME: respect compression
	filehead_ref := base.BlockRef{
		Oid:  ref.Oid,
		Type: ref.Type,
	}

	filehead_bin, err := BlockLoadDecrypt(bs, filehead_ref)
	if err != nil {
		return wire.FileHead{}, fmt.Errorf("failed loading FileHead [%w]", err)
	}

	filehead, err := wire.FileHeadUnmarshal(filehead_bin)
	if err != nil {
		return wire.FileHead{}, fmt.Errorf("failed unmarshalling FileHead [%w]", err)
	}

	return filehead, err
}

// FIXME: use LoadFileHead
func LoadFileControl(bs base.IBlockStore, ref base.BlockRef) (wire.FileControl, error) {
	// load the index block -- strip off the, that's later used used to decrypt the FileControl block
	filehead_ref := base.BlockRef{
		Oid:  ref.Oid,
		Type: ref.Type,
	}

	filehead_bin, err := BlockLoadDecrypt(bs, filehead_ref)
	if err != nil {
		return wire.FileControl{}, fmt.Errorf("failed loading FileHead [%w]", err)
	}

	filehead, err := wire.FileHeadUnmarshal(filehead_bin)
	if err != nil {
		return wire.FileControl{}, fmt.Errorf("failed unmarshalling FileHead [%w]", err)
	}

	return DecryptFileControl(ref.Cipher, ref.Key, filehead.Private)
}
