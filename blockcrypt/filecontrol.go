package blockcrypt

import (
	"fmt"

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

	fctrl := wire.FileControl{}
	fctrl.Unmarshal(fctrl_bin)

	return fctrl, nil
}
