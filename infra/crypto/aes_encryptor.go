package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrDecryptionFailed = errors.New("decryption failed")
)

type EncryptedData struct {
	Nonce      []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
}

type AESEncryptor struct {
	keyManager *KeyManager
}

func NewAESEncryptor(keyManager *KeyManager) *AESEncryptor {
	return &AESEncryptor{
		keyManager: keyManager,
	}
}

func (e *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	key, err := e.keyManager.LoadKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	encrypted := EncryptedData{
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}

	return json.Marshal(encrypted)
}

func (e *AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	key, err := e.keyManager.LoadKey()
	if err != nil {
		return nil, err
	}

	var encrypted EncryptedData
	if err := json.Unmarshal(data, &encrypted); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, encrypted.Nonce, encrypted.Ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

func (e *AESEncryptor) InitializeKey() error {
	return e.keyManager.InitializeKey()
}

func (e *AESEncryptor) KeyExists() bool {
	return e.keyManager.KeyExists()
}
