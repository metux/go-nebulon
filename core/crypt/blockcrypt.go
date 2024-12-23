package crypt

import (
	"fmt"
	"time"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/util"
)

var (
	TraceEncrypt   = false
	ErrNoCryptoKey = fmt.Errorf("encrypted block: no key")
)

func BlockDecrypt(cipher base.CipherType, key []byte, data []byte) ([]byte, error) {
	if TraceEncrypt {
		defer util.TimeTrack(time.Now(), "BlockDecrypt ("+cipher.String()+")")
	}

	switch cipher {
	case base.CipherType_None:
		return data, nil
	case base.CipherType_AES_CBC:
		return AES_Decrypt(data, key)
	case base.CipherType_AES_CBC_ZSTD:
		return AES_ZSTD_Decrypt(data, key)
	default:
		return nil, fmt.Errorf("unsupported cipher type: %s", cipher)
	}
}

func BlockEncrypt(cipher base.CipherType, data []byte) ([]byte, []byte, base.CipherType, error) {
	if TraceEncrypt {
		defer util.TimeTrack(time.Now(), fmt.Sprintf("BlockEncrypt (%s) for %d bytes", cipher.String(), len(data)))
	}

	switch cipher {
	case base.CipherType_None:
		return nil, data, base.CipherType_None, nil
	case base.CipherType_AES_CBC:
		return AES_Encrypt(data)
	case base.CipherType_AES_CBC_ZSTD:
		return AES_ZSTD_Encrypt(data)
	default:
		return nil, nil, cipher, fmt.Errorf("unsupported cipher type: %s", cipher)
	}
}

func BlockLoadDecrypt(bs base.IBlockStore, ref base.BlockRef) ([]byte, error) {
	encrypted, err := bs.GetBlock(ref)
	if err != nil {
		return encrypted, err
	}

	return BlockDecrypt(ref.Cipher, ref.Key, encrypted)
}
