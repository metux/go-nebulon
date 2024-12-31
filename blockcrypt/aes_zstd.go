package blockcrypt

import (
	"fmt"

	"github.com/klauspost/compress/zstd"
	"github.com/metux/go-nebulon/wire"
)

var (
	zipWriter, zipErrW = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	zipReader, zipErrR = zstd.NewReader(nil)
)

func AES_ZSTD_Encrypt(data []byte) ([]byte, []byte, wire.CipherType, error) {
	compressed := zipWriter.EncodeAll(data, make([]byte, 0, len(data)))

	if len(compressed) >= len(data) {
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

	decoded, err := zipReader.DecodeAll(decrypted, make([]byte, 0, len(data)))
	if err != nil {
		return decoded, fmt.Errorf("AES_ZSTD_Decrypt() failed decompressing [%w]", err)
	}
	return decoded, nil
}
