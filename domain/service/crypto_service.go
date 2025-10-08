package service

type CryptoService interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
	InitializeKey() error
	KeyExists() bool
}
