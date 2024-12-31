package blockcrypt

import (
	"fmt"

	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

func AES_ZSTD_Encrypt(data []byte) ([]byte, []byte, wire.CipherType, error) {
	compressed := util.ZipEncode(data)

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

	decoded, err := util.ZipDecode(decrypted)
	if err != nil {
		return decoded, fmt.Errorf("AES_ZSTD_Decrypt() failed decompressing [%w]", err)
	}
	return decoded, nil
}
