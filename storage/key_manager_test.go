package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyManager_InitializeKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(string)
		hasErr bool
	}{
		{
			name:   "succeed: create new key",
			setup:  func(dir string) {},
			hasErr: false,
		},
		{
			name: "succeed: create key with non-existent directory",
			setup: func(dir string) {
				os.RemoveAll(dir)
			},
			hasErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			test.setup(tmpDir)

			km := NewKeyManager(tmpDir)
			err := km.InitializeKey()

			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, km.KeyExists())

				key, err := os.ReadFile(filepath.Join(tmpDir, KeyFileName))
				assert.NoError(t, err)
				assert.Equal(t, KeySize, len(key))
			}
		})
	}
}

func TestKeyManager_KeyExists(t *testing.T) {
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
			result := km.KeyExists()
			assert.Equal(t, test.want, result)
		})
	}
}

func TestKeyManager_LoadKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(string)
		hasErr bool
		errMsg string
	}{
		{
			name: "succeed: load existing key",
			setup: func(dir string) {
				km := NewKeyManager(dir)
				km.InitializeKey()
			},
			hasErr: false,
		},
		{
			name:   "failed: key not found",
			setup:  func(dir string) {},
			hasErr: true,
			errMsg: "encryption key not found",
		},
		{
			name: "failed: invalid key size",
			setup: func(dir string) {
				os.MkdirAll(dir, DirPermission)
				os.WriteFile(filepath.Join(dir, KeyFileName), []byte("invalid"), KeyPermission)
			},
			hasErr: true,
			errMsg: "invalid key size",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			test.setup(tmpDir)

			km := NewKeyManager(tmpDir)
			key, err := km.LoadKey()

			if test.hasErr {
				assert.Error(t, err)
				if test.errMsg != "" {
					assert.Contains(t, err.Error(), test.errMsg)
				}
				assert.Nil(t, key)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, key)
				assert.Equal(t, KeySize, len(key))
			}
		})
	}
}
