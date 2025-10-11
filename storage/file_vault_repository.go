package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ritarock/passvault/domain"
)

const (
	VaultFileName   = "vault.json.enc"
	VaultPermission = 0600
	DirPermission   = 0700
)

var (
	ErrVaultNotFound = errors.New("vault file not found")
)

type FileVaultRepository struct {
	vaultPath string
	cryptoSvc domain.CryptoService
}

func NewFileVaultRepository(baseDir string, cryptoSvc domain.CryptoService) *FileVaultRepository {
	return &FileVaultRepository{
		vaultPath: filepath.Join(baseDir, VaultFileName),
		cryptoSvc: cryptoSvc,
	}
}

func (r *FileVaultRepository) Exists() bool {
	_, err := os.Stat(r.vaultPath)
	return err == nil
}

func (r *FileVaultRepository) Load() (*domain.Vault, error) {
	if !r.Exists() {
		return nil, ErrVaultNotFound
	}

	encryptedData, err := os.ReadFile(r.vaultPath)
	if err != nil {
		return nil, err
	}

	decryptedData, err := r.cryptoSvc.Decrypt(encryptedData)
	if err != nil {
		return nil, err
	}

	var vault domain.Vault
	if err := json.Unmarshal(decryptedData, &vault); err != nil {
		return nil, err
	}

	return &vault, nil
}

func (r *FileVaultRepository) Save(vault *domain.Vault) error {
	dir := filepath.Dir(r.vaultPath)
	if err := os.MkdirAll(dir, DirPermission); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return err
	}

	encryptedData, err := r.cryptoSvc.Encrypt(jsonData)
	if err != nil {
		return err
	}

	return os.WriteFile(r.vaultPath, encryptedData, VaultPermission)
}
