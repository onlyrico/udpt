// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[aes_cipher.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const errAESKeySize = "AES-256 key must be 32 bytes long"

// aesCipher implements the SymmetricCipher interface that encrypts and
// decrypts plaintext using the AES-256 symmetric cipher algorithm.
type aesCipher struct {
	cryptoKey []byte
	gcm       cipher.AEAD
} //                                                                   aesCipher

// ValidateKey checks if 'key' is acceptable for use with the cipher.
// For example it must be of the right size.
//
// For AES-256, the key must be exactly 32 bytes long.
//
func (ob *aesCipher) ValidateKey(key []byte) error {
	if len(key) != 32 {
		return makeError(0xE42FDB, errAESKeySize)
	}
	return nil
} //                                                                 ValidateKey

// SetKey initializes the cipher with the specified secret key.
//
// If the cipher is already initialized with the given key, does nothing.
// The same key is used for encryption and decryption.
//
func (ob *aesCipher) SetKey(key []byte) error {
	err := ob.ValidateKey(key)
	if err != nil {
		return makeError(0xE32BD3, err)
	}
	if bytes.Equal(ob.cryptoKey, key) {
		return nil
	}
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return err
	}
	ob.gcm = gcm
	ob.cryptoKey = key
	return nil
} //                                                                      SetKey

// Encrypt encrypts plaintext using the key given to SetKey and
// returns the encrypted ciphertext, using AES-256 symmetric cipher.
//
// You need to call SetKey at least once before you call Encrypt.
//
func (ob *aesCipher) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	err = ob.ValidateKey(ob.cryptoKey)
	if err != nil {
		return nil, makeError(0xE64A2E, err)
	}
	// nonce is a byte array filled with cryptographically secure random bytes
	n := ob.gcm.NonceSize() // = gcmStandardNonceSize = 12 bytes
	nonce := make([]byte, n)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	ciphertext = ob.gcm.Seal(
		nonce,     // dst
		nonce,     // nonce
		plaintext, // plaintext
		nil,       // additionalData
	)
	return ciphertext, nil
} //                                                                     Encrypt

// Decrypt decrypts ciphertext using the key given to SetKey and
// returns the decrypted plaintext, using AES-256 symmetric cipher.
//
// You need to call SetKey at least once before you call Decrypt.
//
func (ob *aesCipher) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	err = ob.ValidateKey(ob.cryptoKey)
	if err != nil {
		return nil, makeError(0xE35A87, err)
	}
	n := ob.gcm.NonceSize()
	if len(ciphertext) < n {
		return nil, makeError(0xE5F7E2, "invalid ciphertext")
	}
	nonce := ciphertext[:n]
	ciphertext = ciphertext[n:]
	plaintext, err = ob.gcm.Open(
		nil,        // dst
		nonce,      // nonce
		ciphertext, // ciphertext
		nil,        // additionalData
	)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
} //                                                                     Decrypt

// end
