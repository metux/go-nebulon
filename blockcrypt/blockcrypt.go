package blockcrypt

import (
	"log"
	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
)

func main() {
	key := "12345678901234567890123456789012"
	iv := "1234567890123456"
	plaintext := "abcdefghijklmnopqrstuvwxyzABCDEF"
	fmt.Printf("Result: %v\n", Ase256(plaintext, key, iv, aes.BlockSize))
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func Ase256(data[]byte, key []byte, iv []byte) ([]byte, error){
	bPlaintext := PKCS5Padding(data, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return ciphertext, nil
}

func Ase256Decode(cypted [] byte, key [] byte, iv [] byte) ([]byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	decrypted := make([]byte, len(crypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, crypted)
	return PKCS5UnPadding(decrypted)
}

// dummy
func EncryptBlock(data []byte) ([]byte, []byte, error) {
	key := sha256.Sum256(data)
//	iv := key
	iv := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

	encryped, err := Ase256(data, key, iv))
	return encrypted, sha
}

func DecryptBlock(data []byte, key []byte) []byte {
	iv := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

	return data
}
