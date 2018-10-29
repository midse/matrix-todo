package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateKeyFromPassword(password string, salt *[]byte) ([]byte, []byte, error) {
	var err error
	var generatedSalt []byte

	if salt != nil {
		generatedSalt = *salt
	} else {
		// Salt length should be at least the same size as the output of the hash function (sha256 = 32 bytes)
		// https://crackstation.net/hashing-security.htm#salt
		generatedSalt, err = generateRandomBytes(32)

		if err != nil {
			return nil, nil, fmt.Errorf("Unable to securely generate bytes")
		}
	}

	derivedKey := pbkdf2.Key([]byte(password), generatedSalt, 100000, 32, sha256.New)

	return derivedKey, generatedSalt, nil
}

func encrypt(key []byte, message string) (encmess []byte, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Println(err)
		return
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce, err := generateRandomBytes(12)
	if err != nil {
		logger.Println(err)
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Println(err)
		return
	}

	cipherText := aesgcm.Seal(nil, nonce, plainText, nil)
	cipherText = append(nonce, cipherText...)

	encmess = cipherText
	return
}

func decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText := []byte(securemess)
	if err != nil {
		logger.Println(err)
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Println(err)
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		logger.Println(err)
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Println(err)
		return
	}

	plaintext, err := aesgcm.Open(nil, cipherText[0:12], cipherText[12:], nil)

	decodedmess = string(plaintext)
	return
}
