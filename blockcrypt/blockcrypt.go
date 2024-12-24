package blockcrypt

import (
	"crypto/sha256"
)

// dummy
func EncryptBlock(data [] byte) (byte[], byte[]) {
	sha := sha256.Sum256(data)
	return data, sha
}

func DecryptBlock(data [] byte, key [] byte) [] byte {
	return data
}
