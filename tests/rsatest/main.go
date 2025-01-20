package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/metux/go-nebulon/util"
)

func encrypt(pubkey_fn string, s string) []byte {
	pubkey, err := util.LoadRSAPublicKey(pubkey_fn)
	if err != nil {
		log.Fatalf("error loading key: %s", err)
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey, []byte(s), nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

func decrypt(privkey_fn string, ciphertext []byte) string {
	privateKey, err := util.LoadRSAPrivateKey(privkey_fn)
	if err != nil {
		log.Fatalf("error loading privkey %s\n", err)
	}

	decryptedPlaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return string(decryptedPlaintext)
}

func main() {
	text := "hello world hello world hello world hello world hello world hello world"
	log.Printf("text len %d\n", len(text))
	pubkeyfile := "/home/nekrad/.ssh/id_rsa.pub.pem"
	privkeyfile := "/home/nekrad/.ssh/id_rsa"

	encrypted := encrypt(pubkeyfile, text)
	log.Printf("encrypted size %d\n", len(encrypted))
	b64 := base64.StdEncoding.EncodeToString(encrypted)

	log.Printf("Encrypted: %s\n", b64)
	log.Printf("Decrypted: %s\n", decrypt(privkeyfile, encrypted))
}
