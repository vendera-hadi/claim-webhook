package ecies

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ParseECPrivateKeyPEM(keyBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil

}

func ParseECPublicKeyPEM(keyBytes []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected public key type")
	}

	return ecdsaPublicKey, nil
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
