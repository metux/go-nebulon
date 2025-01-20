package main

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/core/crypt"
//	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/util"
)


// 2do: need to rewrite it for hybrid encryption
func decryptAnnounce(encoded_frame [] byte, keyfn string) (wire.AnnouncePayload, error) {
	log.Printf("decryptAnnounce: encoded frame size: %d\n", len(encoded_frame))

	frame, err := wire.AnnounceFrameUnmarshal(encoded_frame)
	if err != nil {
		log.Fatalf("decryptAnnounce: unmarshal frame failed %s\n", err)
		return wire.AnnouncePayload{}, nil
	}
	log.Printf("decryptAnnounce: encrypted payload size %d\n", len(frame.Payload))
	log.Printf("decryptAnnounce: encrypted key size %d\n", len(frame.Key))

	log.Printf("DEC: encrypted session key %v\n", frame.Key)

	decrypted_key, err := util.RSADecrypt(keyfn, frame.Key)
	if err != nil {
		log.Fatalf("decryptAnnounce: decrypting payload key failed: %s\n", err)
	}

	log.Printf("decrypted payload size %d\n", len(decrypted_key))
	decrypted_payload, err := crypt.AES_Decrypt(frame.Payload, decrypted_key)
	if err != nil {
		log.Fatalf("decrypting payload with AES failed\n")
		return wire.AnnouncePayload{}, nil
	}

	log.Printf("decrypted payload len %d\n", len(decrypted_payload))

	pl, err := wire.AnnouncePayloadUnmarshal(decrypted_payload)
	if err != nil {
		log.Fatalf("failed unmarshalling payload: %s\n", err)
	}

	log.Printf("payload decrypt done\n")
	return pl, err
}

func encryptAnnounce(ref wire.BlockRef, keyfn string) ([]byte, error) {
	now := time.Now()
	payload := wire.AnnouncePayload{
		Seconds: now.Unix(),
		Nanos:   now.UnixNano(),
		Ref:     &ref,
	}

	payload_binary, err := payload.Marshal()
	if err != nil {
		return nil, err
	}

	log.Printf("ENC: encoded payload size %d\n", len(payload_binary))

	payload_key, payload_encrypted, _, err := crypt.BlockEncrypt(wire.CipherType_AES_CBC, payload_binary)
	if err != nil {
		log.Fatalf("BlockEncrypt failed: %s\n", err)
	}

	log.Printf("payload key: %v\n", payload_key)

	encrypted_key, err := util.RSAEncrypt(keyfn, payload_key)
	if err != nil {
		log.Fatalf("ENC: frame encrypt error %s\n", err)
		return nil, err
	}

	log.Printf("ENC: encrypted key size %d\n", len(encrypted_key))
	log.Printf("ENC: encrypted payload key %v\n", encrypted_key)

	frame := wire.AnnounceFrame{
		Cipher:  wire.AsymCipherType_RSA,
		Payload: payload_encrypted,
		Key: encrypted_key,
	}

	encoded_frame, err := frame.Marshal()
	if err != nil {
		log.Printf("marshal frame failed\n")
		return nil, err
	}

	log.Printf("encoded frame size %d\n", len(encoded_frame))
	return encoded_frame, nil
}

func main() {
	text := "o world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello worldworld hello worldworld hello worldworld hello world"
	log.Printf("text len %d\n", len(text))
	log.Printf("text data len %d\n", len([]byte(text)))
	pubkeyfile := "/home/nekrad/.ssh/id_rsa.pub.pem"
	privkeyfile := "/home/nekrad/.ssh/id_rsa"

	encrypted, err := util.RSAEncrypt(pubkeyfile, []byte(text))
	if err != nil {
		log.Fatalf("rsa encrypt error %s\n", err)
	}
	log.Printf("encrypted size %d\n", len(encrypted))
	log.Printf("Encrypted: %s\n", base64.StdEncoding.EncodeToString(encrypted))

	decrypted, err := util.RSADecrypt(privkeyfile, encrypted)
	if err != nil {
		log.Fatalf("rsa decrypt failed: %s\n", err)
	}
	decrypted_str := string(decrypted)

	log.Printf("Decrypted: %s\n", decrypted_str)

	if decrypted_str != text {
		log.Fatalf("original and decrypted text mismatch !\n")
	}

	ref := wire.BlockRef{
		Oid: []byte("hello world huhu"),
		Key: []byte("foo bar"),
	}

	encoded_frame, err := encryptAnnounce(ref, pubkeyfile)
	if err != nil {
		log.Printf("encryptAnnounce() error %s\n", err)
	}

	log.Printf("encrypted frame size %d\n", len(encoded_frame))

	_, err = decryptAnnounce(encoded_frame, privkeyfile)
	if err != nil {
		log.Printf("decryptAnnounce failed %s\n", err)
	}
}
