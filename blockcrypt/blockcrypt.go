package blockcrypt

import (
	"fmt"
	"time"

	"github.com/metux/go-nebulon/util"
	"github.com/metux/go-nebulon/wire"
)

var (
	TraceEncrypt = false
)

func BlockDecrypt(cipher wire.CipherType, key []byte, data []byte) ([]byte, error) {
	if TraceEncrypt {
		defer util.TimeTrack(time.Now(), "BlockDecrypt ("+cipher.String()+")")
	}

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

func BlockEncrypt(cipher wire.CipherType, data []byte) ([]byte, []byte, wire.CipherType, error) {
	if TraceEncrypt {
		defer util.TimeTrack(time.Now(), fmt.Sprintf("BlockEncrypt (%s) for %d bytes", cipher.String(), len(data)))
	}

	switch cipher {
	case wire.CipherType_None:
		return []byte{}, data, wire.CipherType_None, nil
	case wire.CipherType_AES_CBC:
		return AES_Encrypt(data)
	case wire.CipherType_AES_CBC_ZSTD:
		return AES_ZSTD_Encrypt(data)
	default:
		return []byte{}, []byte{}, cipher, fmt.Errorf("unsupported cipher type: %s", cipher)
	}
}
