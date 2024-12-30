package blockcrypt

import (
	"fmt"

	"github.com/klauspost/compress/zstd"
)

func AES_ZSTD_Encrypt(data []byte) ([]byte, []byte, error) {
	writer, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("AES_ZSTD_Encrypt() failed creating writer [%w]", err)
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
		return []byte{}, fmt.Errorf("AES_ZSTD_Decrypt() failed creating reader [%w]", err)
	}
	decoded, err := reader.DecodeAll(decrypted, make([]byte, 0, len(data)))
	if err != nil {
		return decoded, fmt.Errorf("AES_ZSTD_Decrypt() failed decompressing [%w]", err)
	}
	return decoded, nil
}
