package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	defaultLength = 16
	minLength     = 8
	maxLength     = 64
	lowercase     = "abcdefghijklmnopqrstuvwxyz"
	uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits        = "0123456789"
	symbols       = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

type PasswordGenerator struct{}

type PasswordOptions struct {
	Length           int
	IncludeLowercase bool
	IncludeUppercase bool
	IncludeDigits    bool
	IncludeSymbols   bool
}

func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{
		Length:           defaultLength,
		IncludeLowercase: true,
		IncludeUppercase: true,
		IncludeDigits:    true,
		IncludeSymbols:   false,
	}
}

func (opts PasswordOptions) Validate() error {
	if opts.Length < minLength {
		return fmt.Errorf("password length must be at least %d", minLength)
	}
	if opts.Length > maxLength {
		return fmt.Errorf("password length must be at most %d", maxLength)
	}
	if !opts.IncludeLowercase && !opts.IncludeUppercase && !opts.IncludeDigits && !opts.IncludeSymbols {
		return fmt.Errorf("at least one character type must be selected")
	}
	return nil
}

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

func (pg *PasswordGenerator) GenerateWithOptions(opts PasswordOptions) (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}

	charset := pg.buildCharset(opts)
	if len(charset) == 0 {
		return "", fmt.Errorf("no characters available for password generation")
	}

	password := make([]byte, opts.Length)
	for i := 0; i < opts.Length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

func (pg *PasswordGenerator) buildCharset(opts PasswordOptions) string {
	var charset string

	if opts.IncludeLowercase {
		charset += lowercase
	}
	if opts.IncludeUppercase {
		charset += uppercase
	}
	if opts.IncludeDigits {
		charset += digits
	}
	if opts.IncludeSymbols {
		charset += symbols
	}

	return charset
}
