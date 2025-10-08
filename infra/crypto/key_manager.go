package crypto

import (
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
)

const (
	KeySize       = 32 // AES-256
	KeyFileName   = "key.bin"
	KeyPermission = 0600
	DirPermission = 0700
)

var (
	ErrKeyNotFound = errors.New("encryption key not found")
)

type KeyManager struct {
	keyPath string
}

func NewKeyManager(baseDir string) *KeyManager {
	return &KeyManager{
		keyPath: filepath.Join(baseDir, KeyFileName),
	}
}

func (km *KeyManager) InitializeKey() error {
	dir := filepath.Dir(km.keyPath)
	if err := os.MkdirAll(dir, DirPermission); err != nil {
		return err
	}

	key := make([]byte, KeySize)
	if _, err := rand.Read(key); err != nil {
		return err
	}

	return os.WriteFile(km.keyPath, key, KeyPermission)
}

func (km *KeyManager) KeyExists() bool {
	_, err := os.Stat(km.keyPath)
	return err == nil
}

func (km *KeyManager) LoadKey() ([]byte, error) {
	if !km.KeyExists() {
		return nil, ErrKeyNotFound
	}
	key, err := os.ReadFile(km.keyPath)
	if err != nil {
		return nil, err
	}

	if len(key) != KeySize {
		return nil, errors.New("invalid key size")
	}

	return key, nil
}
