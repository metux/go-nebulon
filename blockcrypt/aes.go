package blockcrypt

import (
	"crypto/aes"
	"crypto/sha256"
	"fmt"

	"github.com/metux/go-nebulon/util"
)

func ivFromKey(key []byte) []byte {
	hashed := sha256.Sum256(key)
	return hashed[:aes.BlockSize]
}

func AES_Encrypt(data []byte) ([]byte, []byte, error) {
	h := sha256.Sum256(data)
	key := h[:]
	iv := ivFromKey(key)

	encrypted, err := util.AES256Encrypt(data, key, iv)
	if err != nil {
		return key, encrypted, fmt.Errorf("EncryptBlock: encrypting block failed [%w]", err)
	}

	return key, encrypted, nil
}

func AES_Decrypt(data []byte, key []byte) ([]byte, error) {
	return util.AES256Decrypt(data, key, ivFromKey(key))
}
