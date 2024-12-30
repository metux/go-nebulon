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
		return AES_Decrypt(data, key)
	case wire.CipherType_AES_CBC_ZSTD:
		return AES_ZSTD_Decrypt(data, key)
	default:
		return []byte{}, fmt.Errorf("unsupported cipher type: %s", cipher)
	}
}

func BlockEncrypt(cipher wire.CipherType, data []byte) ([]byte, []byte, error) {
	switch cipher {
	case wire.CipherType_None:
		return []byte{}, data, nil
	case wire.CipherType_AES_CBC:
		return AES_Encrypt(data)
	case wire.CipherType_AES_CBC_ZSTD:
		return AES_ZSTD_Encrypt(data)
	default:
		return []byte{}, []byte{}, fmt.Errorf("unsupported cipher type: %s", cipher)
	}
}
