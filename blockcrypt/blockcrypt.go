package blockcrypt

import (
	"fmt"

	"github.com/metux/go-nebulon/wire"
)

func BlockDecrypt(cipher wire.CipherType, key []byte, data []byte) ([]byte, error) {
	switch cipher {
	case wire.CipherType_None:
		return data, nil
	case wire.CipherType_AES_CBC:
		return AESDecryptBlock(data, key)
	case wire.CipherType_AES_CBC_ZSTD:
		return AES_ZSTD_Decrypt(data, key)
	default:
		return []byte{}, fmt.Errorf("unsupported cipher type: %s\n", cipher)
	}
}

func BlockEncrypt(cipher wire.CipherType, data []byte) ([]byte, []byte, error) {
	switch cipher {
	case wire.CipherType_None:
		return []byte{}, data, nil
	case wire.CipherType_AES_CBC:
		encrypted, key, err := AESEncryptBlock(data)
		return key, encrypted, err
	case wire.CipherType_AES_CBC_ZSTD:
		encrypted, key, err := AES_ZSTD_Encrypt(data)
		return key, encrypted, err
	default:
		return []byte{}, []byte{}, fmt.Errorf("storeDataBlock(): unsupported cipher %s\n", cipher)
	}
}
