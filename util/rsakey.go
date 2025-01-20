package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadRSAPublicKey(fn string) (*rsa.PublicKey, error) {
	encoded, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(encoded)
	if block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("wrong block type for public key: %s", block.Type)
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func LoadRSAPrivateKey(fn string) (*rsa.PrivateKey, error) {
	encoded, _ := os.ReadFile(fn)
	block, _ := pem.Decode(encoded)
	if block.Type != "PRIVATE KEY" && block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("wrong block type for private key: %s", block.Type)
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
