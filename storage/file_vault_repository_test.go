package storage

import (
	"errors"
	"testing"

	"github.com/ritarock/passvault/domain"
	"github.com/stretchr/testify/assert"
)

type mockCryptoService struct {
	encryptFunc func([]byte) ([]byte, error)
	decryptFunc func([]byte) ([]byte, error)
}

func (m *mockCryptoService) Encrypt(data []byte) ([]byte, error) {
	if m.encryptFunc != nil {
		return m.encryptFunc(data)
	}
	return data, nil
}

func (m *mockCryptoService) Decrypt(data []byte) ([]byte, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(data)
	}
	return data, nil
}

func (m *mockCryptoService) InitializeKey() error {
	return nil
}

func (m *mockCryptoService) KeyExists() bool {
	return true
}

func TestFileVaultRepository_Save(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func() domain.CryptoService
		vault  *domain.Vault
		hasErr bool
	}{
		{
			name: "succeed: save vault",
			setup: func() domain.CryptoService {
				return &mockCryptoService{
					encryptFunc: func(data []byte) ([]byte, error) {
						return data, nil
					},
				}
			},
			vault:  domain.NewVault(),
			hasErr: false,
		},
		{
			name: "failed: encryption error",
			setup: func() domain.CryptoService {
				return &mockCryptoService{
					encryptFunc: func(data []byte) ([]byte, error) {
						return nil, errors.New("encryption error")
					},
				}
			},
			vault:  domain.NewVault(),
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			cryptoSvc := test.setup()
			repo := NewFileVaultRepository(tmpDir, cryptoSvc)

			err := repo.Save(test.vault)

			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, repo.Exists())
			}
		})
	}
}

func TestFileVaultRepository_Load(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(string) domain.CryptoService
		hasErr bool
		errMsg string
	}{
		{
			name: "succeed: load vault",
			setup: func(tmpDir string) domain.CryptoService {
				cryptoSvc := &mockCryptoService{}
				repo := NewFileVaultRepository(tmpDir, cryptoSvc)
				vault := domain.NewVault()
				repo.Save(vault)
				return cryptoSvc
			},
			hasErr: false,
		},
		{
			name: "failed: vault not found",
			setup: func(tmpDir string) domain.CryptoService {
				return &mockCryptoService{}
			},
			hasErr: true,
			errMsg: "vault file not found",
		},
		{
			name: "failed: decryption error",
			setup: func(tmpDir string) domain.CryptoService {
				cryptoSvc := &mockCryptoService{}
				repo := NewFileVaultRepository(tmpDir, cryptoSvc)
				vault := domain.NewVault()
				repo.Save(vault)

				return &mockCryptoService{
					decryptFunc: func(data []byte) ([]byte, error) {
						return nil, errors.New("decryption error")
					},
				}
			},
			hasErr: true,
			errMsg: "decryption error",
		},
		{
			name: "failed: invalid JSON",
			setup: func(tmpDir string) domain.CryptoService {
				cryptoSvc := &mockCryptoService{}
				repo := NewFileVaultRepository(tmpDir, cryptoSvc)
				vault := domain.NewVault()
				repo.Save(vault)

				return &mockCryptoService{
					decryptFunc: func(data []byte) ([]byte, error) {
						return []byte("invalid json"), nil
					},
				}
			},
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			cryptoSvc := test.setup(tmpDir)
			repo := NewFileVaultRepository(tmpDir, cryptoSvc)

			vault, err := repo.Load()

			if test.hasErr {
				assert.Error(t, err)
				if test.errMsg != "" {
					assert.Contains(t, err.Error(), test.errMsg)
				}
				assert.Nil(t, vault)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, vault)
			}
		})
	}
}

func TestFileVaultRepository_Exists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		setup func(string, domain.CryptoService)
		want  bool
	}{
		{
			name: "succeed: vault exists",
			setup: func(tmpDir string, cryptoSvc domain.CryptoService) {
				repo := NewFileVaultRepository(tmpDir, cryptoSvc)
				vault := domain.NewVault()
				repo.Save(vault)
			},
			want: true,
		},
		{
			name:  "succeed: vault does not exist",
			setup: func(tmpDir string, cryptoSvc domain.CryptoService) {},
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			cryptoSvc := &mockCryptoService{}
			test.setup(tmpDir, cryptoSvc)

			repo := NewFileVaultRepository(tmpDir, cryptoSvc)
			result := repo.Exists()
			assert.Equal(t, test.want, result)
		})
	}
}

func TestFileVaultRepository_SaveAndLoad(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		entry domain.Entry
	}{
		{
			name: "succeed: save and load vault with entry",
			entry: domain.Entry{
				ID:       "test-id-1",
				Title:    "test title",
				Password: "test password",
				URL:      "test url",
				Notes:    "test notes",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()

			km := NewKeyManager(tmpDir)
			km.InitializeKey()
			encryptor := NewAESEncryptor(km)
			repo := NewFileVaultRepository(tmpDir, encryptor)

			vault := domain.NewVault()
			vault.CreateEntry(test.entry)

			err := repo.Save(vault)
			assert.NoError(t, err)

			loadedVault, err := repo.Load()
			assert.NoError(t, err)
			assert.NotNil(t, loadedVault)
			assert.Equal(t, vault.Version, loadedVault.Version)
			assert.Equal(t, len(vault.Entries), len(loadedVault.Entries))

			loadedEntry, err := loadedVault.GetEntry(test.entry.ID)
			assert.NoError(t, err)
			assert.Equal(t, test.entry.Title, loadedEntry.Title)
			assert.Equal(t, test.entry.Password, loadedEntry.Password)
			assert.Equal(t, test.entry.URL, loadedEntry.URL)
			assert.Equal(t, test.entry.Notes, loadedEntry.Notes)
		})
	}
}
