package blockcrypt

import (
	"crypto/aes"
	"crypto/sha256"
	"log"

	"github.com/metux/go-nebulon/util"
	"github.com/klauspost/compress/zstd"
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
		log.Printf("EncryptBlock: encrypting block failed: %s\n", err)
	}

	return key, encrypted, nil
}

func AES_Decrypt(data []byte, key []byte) ([]byte, error) {
	return util.AES256Decrypt(data, key, ivFromKey(key))
}

func AES_ZSTD_Encrypt(data []byte) ([]byte, []byte, error) {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		log.Printf("AES_ZSTD_Encrypt() failed creating writer: %s\n", err)
		return []byte{}, []byte{}, err
	}
	compressed := writer.EncodeAll(data, make([]byte, 0, len(data)))

	return AES_Encrypt(compressed)
}

func AES_ZSTD_Decrypt(data []byte, key []byte) ([]byte, error) {
	decrypted, err := AES_Decrypt(data, key)
	if err != nil {
		return []byte{}, err
	}

	reader, err := zstd.NewReader(nil)
	if err != nil {
		log.Printf("AES_ZSTD_Decrypt() failed creating reader: %s\n", err)
		return []byte{}, err
	}
	decoded, err := reader.DecodeAll(decrypted, make([]byte, 0, len(data)))
	if err != nil {
		log.Printf("AES_ZSTD_Decrypt() failed decompressing %s\n", err)
		return decoded, err
	}
	return decoded, nil
}
