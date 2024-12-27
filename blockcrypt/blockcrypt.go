package blockcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"log"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

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

func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func IvFromKey(key []byte) []byte {
	hashed := sha256.Sum256(key)
	return hashed[:aes.BlockSize]
}

func EncryptBlock(data []byte) ([]byte, []byte, error) {
	key := Hash(data)
	iv := IvFromKey(key)

	encrypted, err := AES256Encrypt(data, key, iv)
	if err != nil {
		log.Printf("EncryptBlock: encrypting block failed: %s\n", err)
	}

	return encrypted, key, nil
}

func DecryptBlock(data []byte, key []byte) []byte {
	iv := IvFromKey(key)

	data, err := AES256Decrypt(data, key, iv)
	if err != nil {
		panic(err)
	}
	return data
}
