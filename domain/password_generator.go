package domain

import (
	"crypto/rand"
	"math/big"
)

const (
	defaultLength = 16
	lowercase     = "abcdefghijklmnopqrstuvwxyz"
	uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits        = "0123456789"
	symbols       = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

type PasswordGenerator struct{}

func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

// Generate generates a random password with the specified length
// Default length is 16 characters
// The password will contain uppercase, lowercase, digits, and symbols
func (pg *PasswordGenerator) Generate(length int) (string, error) {
	if length <= 0 {
		length = defaultLength
	}

	charset := lowercase + uppercase + digits + symbols
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}
