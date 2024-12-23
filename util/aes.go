package util

import (
	"crypto/aes"
	"crypto/cipher"
)

func AES256Encrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	bPlaintext := PKCS5Padding(data, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return ciphertext, nil
}

func AES256Decrypt(crypted []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	decrypted := make([]byte, len(crypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, crypted)
	return PKCS5UnPadding(decrypted), nil
}
