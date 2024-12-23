package crypt

import (
	"crypto/aes"
	"fmt"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/util"
)

func ivFromKey(key []byte) []byte {
	hashed := util.ContentKey(key)
	return hashed[:aes.BlockSize]
}

func AES_Encrypt(data []byte) ([]byte, []byte, base.CipherType, error) {
	key := util.ContentKey(data)
	iv := ivFromKey(key)

	encrypted, err := util.AES256Encrypt(data, key, iv)
	if err != nil {
		return key, encrypted, base.CipherType_None, fmt.Errorf("EncryptBlock: encrypting block failed [%w]", err)
	}

	return key, encrypted, base.CipherType_AES_CBC, nil
}

func AES_Decrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, ErrNoCryptoKey
	}
	return util.AES256Decrypt(data, key, ivFromKey(key))
}
