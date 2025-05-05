package infra

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ritarock/pa55vault/entity"
)

type Store struct {
	filePath string
}

func NewStore(filePath string) *Store {
	return &Store{
		filePath: filePath,
	}
}

func (s *Store) Read() ([]entity.Vault, error) {
	f, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var vaults []entity.Vault
	if err := json.Unmarshal(data, &vaults); err != nil {
		return nil, err
	}

	return vaults, nil
}

func (s *Store) Write(vaults []entity.Vault) error {
	data, err := json.MarshalIndent(vaults, "", " ")
	if err != nil {
		return err
	}

	f, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
