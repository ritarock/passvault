package infra

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_Read(t *testing.T) {
	t.Parallel()
	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tmpFile, err := os.CreateTemp("./", "test_read.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []testData{
		{ID: 1, Name: "testA"},
		{ID: 2, Name: "testB"},
	}

	data, err := json.Marshal(content)
	assert.NoError(t, err)

	err = os.WriteFile(tmpFile.Name(), data, 0644)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		store Store[testData]
		want  []testData
	}{
		{
			name: "pass",
			store: Store[testData]{
				filepath: tmpFile.Name(),
			},
			want: []testData{
				{ID: 1, Name: "testA"},
				{ID: 2, Name: "testB"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.store.Read()
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestStore_Write(t *testing.T) {
	t.Parallel()
	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tmpFile, err := os.CreateTemp("./", "test_write.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	assert.NoError(t, err)

	tests := []struct {
		name  string
		store Store[testData]
		data  []testData
	}{
		{
			name: "pass",
			store: Store[testData]{
				filepath: tmpFile.Name(),
			},
			data: []testData{
				{ID: 1, Name: "testA"},
				{ID: 2, Name: "testB"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.store.Write(test.data)
			assert.NoError(t, err)
		})
	}
}
