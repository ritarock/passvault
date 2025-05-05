package entity

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func (v *Vault) Encrypt() error {
	key := createKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	code := pkcs7Pad([]byte(v.Code), aes.BlockSize)
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	ciphertext := make([]byte, len(code))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, code)

	v.Iv = hex.EncodeToString(iv)
	v.Code = base64.StdEncoding.EncodeToString(ciphertext)

	return nil
}

func (v *Vault) Decrypt() (string, error) {
	key := createKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(v.Code)
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(v.Iv)
	if err != nil {
		return "", err
	}

	if len(iv) != aes.BlockSize {
		return "", errors.New("failed decrypt")
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("failed decrypt")
	}

	decrypted := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, ciphertext)

	return string(pkcs7Unpad(decrypted)), nil
}

func createKey() []byte {
	return []byte("12345678901234567890123456789012") // 32byte (AES-256)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func pkcs7Unpad(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	padLen := int(data[len(data)-1])
	if padLen > len(data) {
		return nil
	}
	return data[:len(data)-padLen]
}
