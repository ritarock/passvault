package repository

import "github.com/ritarock/passvault/domain/entity"

type VaultRepository interface {
	Load() (*entity.Vault, error)
	Save(vault *entity.Vault) error
	Exists() bool
}
