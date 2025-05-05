package infra

import (
	"os"
	"testing"

	"github.com/ritarock/pa55vault/entity"
	"github.com/stretchr/testify/assert"
)

func TestStore_Read(t *testing.T) {
	tmpFile, err := os.CreateTemp("./", "test_read.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := `[{"title": "title", "url": "url", "code": "code"}]`
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		store Store
		want  []entity.Vault
	}{
		{
			name: "pass",
			store: Store{
				filePath: tmpFile.Name(),
			},
			want: []entity.Vault{
				{Title: "title", Url: "url", Code: "code"},
			},
		},
	}

	for _, test := range tests {
		got, err := test.store.Read()
		assert.NoError(t, err)
		assert.Equal(t, test.want, got)
	}
}

func TestStore_Write(t *testing.T) {
	tmpFile, err := os.CreateTemp("./", "test_write.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	assert.NoError(t, err)

	tests := []struct {
		name   string
		store  Store
		vaults []entity.Vault
	}{
		{
			name: "pass",
			store: Store{
				filePath: tmpFile.Name(),
			},
			vaults: []entity.Vault{
				{Title: "title", Url: "url", Code: "code"},
			},
		},
	}

	for _, test := range tests {
		err := test.store.Write(test.vaults)
		assert.NoError(t, err)
	}
}
