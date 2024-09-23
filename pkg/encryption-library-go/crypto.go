package ecies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"time"

	"golang.org/x/crypto/hkdf"
)

func Encrypt(message []byte, senderPrivateKey *ecdsa.PrivateKey, receiverPublicKey *ecdsa.PublicKey) ([]byte, error) {
	// Generate ephemeral key
	ephemeralKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	ecdhEphemeralKey, err := ephemeralKey.ECDH()
	if err != nil {
		return nil, err
	}

	// Create nonce, 12 bytes for AES GCM
	nonce := make([]byte, 12) // 96 bits for nonce/IV
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	/// Calculate shared secret using scalar multiplication
	ecdhPublicKey, _ := receiverPublicKey.ECDH()
	Z, _ := ecdhEphemeralKey.ECDH(ecdhPublicKey)

	// Derives symmetric encryption key from Z and ephemeralKey.PublicKey
	z := hkdf.New(sha256.New, Z, ephemeralKey.PublicKey.X.Bytes(), ephemeralKey.PublicKey.Y.Bytes())
	kSym := make([]byte, 32)
	_, err = z.Read(kSym)
	if err != nil {
		return nil, err
	}

	// Initialize AES GCM Block cipher
	block, err := aes.NewCipher(kSym[:])
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Get unix timestamp for additional data
	now := time.Now().Unix()
	additionalData := make([]byte, 8)
	binary.LittleEndian.PutUint64(additionalData, uint64(now))

	// Encrypt the message
	ciphertext := aesGCM.Seal(nil, nonce, message, additionalData)

	// Create signature
	hash := sha256.Sum256(append(ciphertext, additionalData...))
	r, s, err := ecdsa.Sign(rand.Reader, senderPrivateKey, hash[:])
	if err != nil {
		return nil, err
	}
	signature := append(r.Bytes(), s.Bytes()...)

	// Construct the ECIES message
	encrypted := append(ephemeralKey.PublicKey.X.Bytes(), ephemeralKey.PublicKey.Y.Bytes()...)
	encrypted = append(encrypted, nonce...)
	encrypted = append(encrypted, signature...)
	encrypted = append(encrypted, ciphertext...)
	encrypted = append(encrypted, additionalData...)

	return encrypted, nil
}

func Decrypt(encrypted []byte, receiverPrivateKey *ecdsa.PrivateKey, senderPublicKey *ecdsa.PublicKey) ([]byte, error) {

	keySize := receiverPrivateKey.Curve.Params().BitSize / 8
	signatureLength := 64

	/// Recosntruct ephemeral public key from encrypted message
	ephemeralPublicKeyX := encrypted[:keySize]
	ephemeralPublicKeyY := encrypted[keySize : 2*keySize]
	ephemeralPublicKey := ecdsa.PublicKey{
		Curve: receiverPrivateKey.Curve,
		X:     new(big.Int).SetBytes(ephemeralPublicKeyX),
		Y:     new(big.Int).SetBytes(ephemeralPublicKeyY),
	}

	ecdhEphemeralPublicKey, _ := ephemeralPublicKey.ECDH()
	ecdhReceiverPrivateKey, _ := receiverPrivateKey.ECDH()

	// Calculate shared secret
	Z, _ := ecdhReceiverPrivateKey.ECDH(ecdhEphemeralPublicKey)

	// Derives symmetric encryption key from Z and ephemeralKey.PublicKey
	z := hkdf.New(sha256.New, Z, ephemeralPublicKeyX, ephemeralPublicKeyY)
	kSym := make([]byte, 32)
	_, err := z.Read(kSym)
	if err != nil {
		return nil, err
	}

	// Create a new AES-GCM cipher block mode
	block, err := aes.NewCipher(kSym[:])
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract the remaining items
	additionalDataSize := 8
	nonceSize := aesGCM.NonceSize()
	nonce := encrypted[2*keySize : 2*keySize+nonceSize]
	signature := encrypted[2*keySize+nonceSize : 2*keySize+nonceSize+signatureLength]
	ciphertext := encrypted[2*keySize+nonceSize+signatureLength : len(encrypted)-additionalDataSize]
	additionalData := encrypted[len(encrypted)-additionalDataSize:]

	// Signature verification
	hashCalc := sha256.Sum256(append(ciphertext, additionalData...))
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	if !ecdsa.Verify(senderPublicKey, hashCalc[:], r, s) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Check for valid timestamp
	unixTime := int64(binary.LittleEndian.Uint64(additionalData))
	delta := Abs(time.Now().Unix() - unixTime)
	if delta > 20 {
		return nil, fmt.Errorf("invalid timestamp")
	}

	// Decrypt the ciphertext using AES-GCM
	message, err := aesGCM.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, err
	}

	return message, nil

}
