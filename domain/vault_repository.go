package domain

type VaultRepository interface {
	Load() (*Vault, error)
	Save(vault *Vault) error
	Exists() bool
}
