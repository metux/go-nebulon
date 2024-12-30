package blockcrypt

import (
	"crypto/aes"
	"crypto/sha256"
	"fmt"

	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

func ivFromKey(key []byte) []byte {
	hashed := sha256.Sum256(key)
	return hashed[:aes.BlockSize]
}

func AES_Encrypt(data []byte) ([]byte, []byte, wire.CipherType, error) {
	h := sha256.Sum256(data)
	key := h[:]
	iv := ivFromKey(key)

	encrypted, err := util.AES256Encrypt(data, key, iv)
	if err != nil {
		return key, encrypted, wire.CipherType_None, fmt.Errorf("EncryptBlock: encrypting block failed [%w]", err)
	}

	return key, encrypted, wire.CipherType_AES_CBC, nil
}

func AES_Decrypt(data []byte, key []byte) ([]byte, error) {
	return util.AES256Decrypt(data, key, ivFromKey(key))
}
