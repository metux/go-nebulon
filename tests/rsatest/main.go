package main

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/util"
)

func encryptAnnounce(ref wire.BlockRef, keyfn string) ([]byte, error) {
	now := time.Now()
	frame := wire.AnnouncePayload{
		Seconds: now.Unix(),
		Nanos:   now.UnixNano(),
		Ref:     &ref,
	}

	encoded_frame, err := frame.Marshal()
	if err != nil {
		log.Printf("frame marshal error: %s\n", err)
		return nil, err
	}

	log.Printf("encoded frame size %d\n", len(encoded_frame))

	encrypted_frame, err := util.RSAEncrypt(keyfn, encoded_frame)
	if err != nil {
		log.Printf("frame encrypt error %s\n", err)
	}

	return encrypted_frame, nil
}

func main() {
	text := "hello world hello world hello world hello world hello world hello world"
	log.Printf("text len %d\n", len(text))
	pubkeyfile := "/home/nekrad/.ssh/id_rsa.pub.pem"
	privkeyfile := "/home/nekrad/.ssh/id_rsa"

	encrypted, err := util.RSAEncrypt(pubkeyfile, []byte(text))
	if err != nil {
		log.Fatalf("rsa encrypt error %s\n")
	}
	log.Printf("encrypted size %d\n", len(encrypted))
	log.Printf("Encrypted: %s\n", base64.StdEncoding.EncodeToString(encrypted))

	decrypted, err := util.RSADecrypt(privkeyfile, encrypted)
	if err != nil {
		log.Fatalf("rsa decrypt failed: %s\n", err)
	}
	log.Printf("Decrypted: %s\n", string(decrypted))

	ref := wire.BlockRef{
		Oid: []byte("hello world huhu"),
		Key: []byte("foo bar"),
	}

	encoded_frame, err := encryptAnnounce(ref, pubkeyfile)
	if err != nil {
		log.Printf("encryptAnnounce() error %s\n", err)
	}

	log.Printf("encrypted frame size %d\n", len(encoded_frame))
}
