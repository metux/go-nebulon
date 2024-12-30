package blockcrypt

import (
	"fmt"

	"github.com/metux/go-nebulon/wire"
)

func EncryptFileControl(content_ref wire.BlockRef, headers map[string]string, cipher wire.CipherType) ([]byte, []byte, error) {

	// fixme: add headers
	fctrl := wire.FileControl{
		Content: &content_ref,
		Headers: headers,
	}
	data, err := fctrl.Marshal()
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("failed marshalling FileControl [%w]", err)
	}

	key, encrypted, err := BlockEncrypt(cipher, data)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("failed encrypting FileControl [%w]", err)
	}

	return key, encrypted, nil
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
