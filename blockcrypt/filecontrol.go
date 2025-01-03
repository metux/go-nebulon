package blockcrypt

import (
	"fmt"

	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/wire"
)

func EncryptFileControl(cipher wire.CipherType, fctrl wire.FileControl) ([]byte, []byte, wire.CipherType, error) {
	data, err := fctrl.Marshal()
	if err != nil {
		return []byte{}, []byte{}, cipher, fmt.Errorf("failed marshalling FileControl [%w]", err)
	}

	key, encrypted, cipher, err := BlockEncrypt(cipher, data)
	if err != nil {
		return []byte{}, []byte{}, cipher, fmt.Errorf("failed encrypting FileControl [%w]", err)
	}

	return key, encrypted, cipher, nil
}

func DecryptFileControl(cipher wire.CipherType, key []byte, private []byte) (wire.FileControl, error) {
	fctrl_bin, err := BlockDecrypt(cipher, key, private)
	if err != nil {
		return wire.FileControl{}, fmt.Errorf("error decrypting FileControl [%w]", err)
	}

	return wire.FileControlUnmarshal(fctrl_bin)
}

func LoadFileControl(bs base.BlockStore, ref wire.BlockRef) (wire.FileControl, error) {
	// load the index block -- strip off the, that's later used used to decrypt the FileControl block
	filehead_ref := wire.BlockRef{
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
