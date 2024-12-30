package blockcrypt

import (
	"fmt"
	"log"

	"github.com/klauspost/compress/zstd"
	"github.com/metux/go-nebulon/wire"
)

func AES_ZSTD_Encrypt(data []byte) ([]byte, []byte, wire.CipherType, error) {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return []byte{}, []byte{}, wire.CipherType_None, fmt.Errorf("AES_ZSTD_Encrypt() failed creating writer [%w]", err)
	}
	compressed := writer.EncodeAll(data, make([]byte, 0, len(data)))
	log.Printf("raw %d compress %d ratio %f\n", len(data), len(compressed), float32(len(compressed))/float32(len(data)))

	if len(compressed) >= len(data) {
		log.Printf("no size improvement - using uncompressed")
		return AES_Encrypt(data)
	}

	key, encrypted, _, err := AES_Encrypt(compressed)
	return key, encrypted, wire.CipherType_AES_CBC_ZSTD, err
}

func AES_ZSTD_Decrypt(data []byte, key []byte) ([]byte, error) {
	decrypted, err := AES_Decrypt(data, key)
	if err != nil {
		return []byte{}, err
	}

	reader, err := zstd.NewReader(nil)
	if err != nil {
		return []byte{}, fmt.Errorf("AES_ZSTD_Decrypt() failed creating reader [%w]", err)
	}
	decoded, err := reader.DecodeAll(decrypted, make([]byte, 0, len(data)))
	if err != nil {
		return decoded, fmt.Errorf("AES_ZSTD_Decrypt() failed decompressing [%w]", err)
	}
	return decoded, nil
}
