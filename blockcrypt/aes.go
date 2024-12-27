package blockcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"log"

	"github.com/klauspost/compress/zstd"
)

func aes256Encrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
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

func aes256Decrypt(crypted []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	decrypted := make([]byte, len(crypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, crypted)
	return PKCS5UnPadding(decrypted), nil
}

func ivFromKey(key []byte) []byte {
	hashed := sha256.Sum256(key)
	return hashed[:aes.BlockSize]
}

func AESEncryptBlock(data []byte) ([]byte, []byte, error) {
	h := sha256.Sum256(data)
	key := h[:]
	iv := ivFromKey(key)

	encrypted, err := aes256Encrypt(data, key, iv)
	if err != nil {
		log.Printf("EncryptBlock: encrypting block failed: %s\n", err)
	}

	return encrypted, key, nil
}

func AESDecryptBlock(data []byte, key []byte) []byte {
	iv := ivFromKey(key)

	data, err := aes256Decrypt(data, key, iv)
	if err != nil {
		panic(err)
	}
	return data
}

func AES_ZSTD_Encrypt(data []byte) ([]byte, []byte, error) {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		log.Printf("AES_ZSTD_Encrypt() failed creating writer: %s\n", err)
		return []byte{}, []byte{}, err
	}
	compressed := writer.EncodeAll(data, make([]byte, 0, len(data)))

	return AESEncryptBlock(compressed)
}

func AES_ZSTD_Decrypt(data []byte, key []byte) []byte {
	decrypted := AESDecryptBlock(data, key)
	reader, err := zstd.NewReader(nil)
	if err != nil {
		log.Printf("AES_ZSTD_Decrypt() failed creating reader: %s\n", err)
		return []byte{}
	}
	decoded, err := reader.DecodeAll(decrypted, make([]byte, 0, len(data)))
	if err != nil {
		log.Printf("AES_ZSTD_Decrypt() failed decompressing %s\n", err)
		return decoded
	}
	return decoded
}
