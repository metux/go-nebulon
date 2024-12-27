package blockcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/metux/go-nebulon/wire"
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

func IvFromKey(key []byte) []byte {
	hashed := sha256.Sum256(key)
	return hashed[:aes.BlockSize]
}

func AESEncryptBlock(data []byte) ([]byte, []byte, error) {
	h := sha256.Sum256(data)
	key := h[:]
	iv := IvFromKey(key)

	encrypted, err := AES256Encrypt(data, key, iv)
	if err != nil {
		log.Printf("EncryptBlock: encrypting block failed: %s\n", err)
	}

	return encrypted, key, nil
}

func AESDecryptBlock(data []byte, key []byte) []byte {
	iv := IvFromKey(key)

	data, err := AES256Decrypt(data, key, iv)
	if err != nil {
		panic(err)
	}
	return data
}

func BlockDecrypt(cipher wire.CipherType, key [] byte, data [] byte) ([]byte, error) {
	switch (cipher) {
		case wire.CipherType_None:
			return data, nil
		case wire.CipherType_AES_CBC:
			return AESDecryptBlock(data, key), nil
		default:
			return []byte{}, fmt.Errorf("unsupported cipher type: %s\n", cipher)
	}
}
