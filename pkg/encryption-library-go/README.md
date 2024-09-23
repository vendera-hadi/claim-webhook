# encryption-library-go
## Go Cryptographic Functions in SATUSEHAT Platform

### Introduction
This module provides a simple and straightforward API for encryption and decryption operations. Developers can integrate SATUSEHAT Cryptographic functions into their Go projects by importing the module and utilizing the **Encrypt** and **Decrypt** functions. The **Encrypt** function takes the sender's private key, the receiver's public key, and plaintext data as input and returns the encrypted data. Similarly, the **Decrypt** function requires the receiver's private key, sender's public key, and encrypted data as input and returns the decrypted data.

It aims to simplify the process of incorporating ECC encryption and decryption capabilities into applications.

### Workflow
The standard based on ECIES (Elliptic Curve Integrated Encryption Scheme) which is a hybrid encryption scheme that combines symmetric encryption, asymmetric encryption, and cryptographic hashing to provide secure communication using Elliptic Curve Cryptography (ECC). The generation of an ECIES message involves several steps:

<ol>
  <li>Key Generation:
    <ul>
      <li>The sender generates a random ephemeral private key, denoted as $kEph$.</li>
      <li>The sender derives the corresponding ephemeral public key, denoted as $KEph$, and where $KEph$ is the base point on the P-256 elliptic curve.</li>
      <li>The sender also obtains the recipient's public key, denoted as $KRcp$, through a pre-established key exchange mechanism.</li>
      <li>The sender computes a shared secret, denoted as $Z$, by performing $ECDH$ function (scalar multiplication) of the recipient's public key $KRcp$ with the sender's ephemeral private key $kEph$. $$Z=ECDH(KRcp,kEph)$$</li>
      <li>The sender derives a symmetric encryption key, denoted as $kSym$, for the symmetric encryption algorithm (e.g., AES-GCM) using HKDF (HMAC Based Key Derivation Function). $$kSym = HKDF(Z,KEph)$$</li>
    </ul>
  </li>
<li>Encryption:
  <ul>
    <li>The sender generate nonce, denoted as $N$</li>
    <li>The sender encrypts the plaintext message, denoted as $M$, using the symmetric encryption algorithm (AES-GCM) with the key $kSym$ and nonce $N$, resulting in the ciphertext, denoted as $C$. $$C=AES_{kSym}(N,M)$$</li>
  </ul>
</li>
<li>Signature Creation:
  <ul>
    <li>The sender calculate hash value from ciphertext $C$, denoted as $h$. $$h=H(C)$$</li>
    <li>The sender generate signature using sender's private key $kSnd$. $$Sig=S(h,kSnd)$$</li>
  </ul>
</li>

<li>Composition: The final ECIES message, denoted as $P$, consists of the following components (ordered):
  $$P=(KEph || N || Sig || C)$$
  <ul>
    <li>The sender's ephemeral public key $KEph$.</li>
    <li>The Nonce used for symmetric encryption $N$</li>
    <li>The ciphertext signature $Sig$</li>
    <li>The ciphertext $C$</li>
  </ul>
</li>  
</ol>


### Installation
To use this module in your Go project, you need to have Go installed. Once you have Go set up, you can install this module using the following command:
```console
go get gitlab.com/dto-moh/satusehat/encryption-library-go
```

### Usage
Import the module in your Go code:
```go
import (
  ss_crypto "gitlab.com/dto-moh/satusehat/encryption-library-go"
)
```

### Encrypt
To encrypt data, use the Encrypt function:
```go
// retrieve sender private key (own private key) from environment variables or from secret manager
senderPrivateKeyPEMBytes := []byte(os.Getenv("PRIVATE_KEY"))
senderPrivateKey, _ := ss_crypto.ParseECPrivateKeyPEM(senderPrivateKeyPEMBytes)

// read receiver public key from public key repository
receiverPublicPEMBytes := []bytes("-----BEGIN PUBLIC KEY----- .... -----END PUBLIC KEY-----")
receiverPublicKey, _ := ss_crypto.ParseECPublicKeyPEM(receiverPublicPEMBytes)

message := "Hello World!"
encrypted, err := ss_crypto.Encrypt([]byte(message), senderPrivateKey, receiverPublicKey)
if err != nil {
   // handle error
}
```

### Decrypt
To decrypt the encrypted data, use the Decrypt function:
```go
// retrieve receiver private key (own private key) from environment variables or from secret manager
receiverPrivateKeyPEMBytes := []byte(os.Getenv("PRIVATE_KEY"))
receiverPrivateKey, _ := ss_crypto.ParseECPrivateKeyPEM(receiverPrivateKeyPEMBytes)

// read sender public key from public key repository
senderPublicPEMBytes := []bytes("-----BEGIN PUBLIC KEY----- .... -----END PUBLIC KEY-----")
senderPublicKey, _ := ss_crypto.ParseECPublicKeyPEM(senderPublicPEMBytes)

plaintext, err := ss_crypto.Decrypt(encrypted, receiverPrivateKey, senderPublicKey)
if err != nil {
   // handle error
}
```

### Generate Key Pairs
To generate ECC Key Pairs using openssl:
```console
openssl ecparam -name prime256v1 -genkey -noout -out private-key.pem
openssl ec -in private-key.pem -pubout -out public-key.pem
```

  
### Reference
- https://asecuritysite.com/ecies/nhsx
- https://cryptobook.nakov.com/asymmetric-key-ciphers/elliptic-curve-cryptography-ecc
- https://blog.cloudflare.com/padding-oracles-and-the-decline-of-cbc-mode-ciphersuites/
- https://crypto.stackexchange.com/questions/202/should-we-mac-then-encrypt-or-encrypt-then-mac
- https://link.springer.com/content/pdf/10.1007/3-540-44448-3_41.pdf
- https://malware.news/t/everyone-loves-curves-but-which-elliptic-curve-is-the-most-popular/17657
- https://github.com/danielhavir/go-hpke
- https://github.com/dajiaji/hpke-js
- https://medium.com/@nimanthaF/authenticated-encryption-with-association-data-aead-ce51db79ddca
- https://en.wikipedia.org/wiki/Integrated_Encryption_Scheme
- https://blog.cloudflare.com/padding-oracles-and-the-decline-of-cbc-mode-ciphersuites/

