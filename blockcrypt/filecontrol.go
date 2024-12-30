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
