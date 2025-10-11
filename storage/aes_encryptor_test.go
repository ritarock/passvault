package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESEncryptor_Encrypt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(string) *KeyManager
		data   []byte
		hasErr bool
	}{
		{
			name: "succeed: encrypt data",
			setup: func(dir string) *KeyManager {
				km := NewKeyManager(dir)
				km.InitializeKey()
				return km
			},
			data:   []byte("test data"),
			hasErr: false,
		},
		{
			name: "failed: key not found",
			setup: func(dir string) *KeyManager {
				return NewKeyManager(dir)
			},
			data:   []byte("test data"),
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			km := test.setup(tmpDir)
			encryptor := NewAESEncryptor(km)

			encrypted, err := encryptor.Encrypt(test.data)

			if test.hasErr {
				assert.Error(t, err)
				assert.Nil(t, encrypted)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, encrypted)
				assert.NotEqual(t, test.data, encrypted)

				var encData EncryptedData
				err := json.Unmarshal(encrypted, &encData)
				assert.NoError(t, err)
				assert.NotEmpty(t, encData.Nonce)
				assert.NotEmpty(t, encData.Ciphertext)
			}
		})
	}
}

func TestAESEncryptor_Decrypt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(string, []byte) ([]byte, error)
		hasErr bool
	}{
		{
			name: "succeed: decrypt data",
			setup: func(dir string, data []byte) ([]byte, error) {
				km := NewKeyManager(dir)
				km.InitializeKey()
				encryptor := NewAESEncryptor(km)
				return encryptor.Encrypt(data)
			},
			hasErr: false,
		},
		{
			name: "failed: invalid encrypted data",
			setup: func(dir string, data []byte) ([]byte, error) {
				km := NewKeyManager(dir)
				km.InitializeKey()
				return []byte("invalid data"), nil
			},
			hasErr: true,
		},
		{
			name: "failed: tampered ciphertext",
			setup: func(dir string, data []byte) ([]byte, error) {
				km := NewKeyManager(dir)
				km.InitializeKey()
				encryptor := NewAESEncryptor(km)
				encrypted, err := encryptor.Encrypt(data)
				if err != nil {
					return nil, err
				}

				var encData EncryptedData
				json.Unmarshal(encrypted, &encData)
				encData.Ciphertext[0] ^= 0xFF
				return json.Marshal(encData)
			},
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			originalData := []byte("test data")

			encrypted, err := test.setup(tmpDir, originalData)
			assert.NoError(t, err)

			km := NewKeyManager(tmpDir)
			encryptor := NewAESEncryptor(km)
			decrypted, err := encryptor.Decrypt(encrypted)

			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, originalData, decrypted)
			}
		})
	}
}

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "succeed: encrypt and decrypt sensitive information",
			data: []byte("sensitive information"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			km := NewKeyManager(tmpDir)
			km.InitializeKey()
			encryptor := NewAESEncryptor(km)

			encrypted, err := encryptor.Encrypt(test.data)
			assert.NoError(t, err)
			assert.NotEqual(t, test.data, encrypted)

			decrypted, err := encryptor.Decrypt(encrypted)
			assert.NoError(t, err)
			assert.Equal(t, test.data, decrypted)
		})
	}
}

func TestAESEncryptor_InitializeKey(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	km := NewKeyManager(tmpDir)
	encryptor := NewAESEncryptor(km)

	err := encryptor.InitializeKey()
	assert.NoError(t, err)
	assert.True(t, encryptor.KeyExists())
}

func TestAESEncryptor_KeyExists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		setup func(string)
		want  bool
	}{
		{
			name: "succeed: key exists",
			setup: func(dir string) {
				km := NewKeyManager(dir)
				km.InitializeKey()
			},
			want: true,
		},
		{
			name:  "succeed: key does not exist",
			setup: func(dir string) {},
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			test.setup(tmpDir)

			km := NewKeyManager(tmpDir)
			encryptor := NewAESEncryptor(km)
			result := encryptor.KeyExists()
			assert.Equal(t, test.want, result)
		})
	}
}
