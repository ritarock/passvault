package core

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Vault struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Code  string `json:"code"`
	Iv    string `json:"iv"`
}

const (
	Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase = "abcdefghijklmnopqrstuvwxyz"
	Number    = "0123456789"
	Symbol    = "!#$%&()-@_<>"
)

const TokenLength = 12

func (v *Vault) GenerateCode() error {
	var token []byte
	for _, charset := range []string{Uppercase, Lowercase, Number, Symbol} {
		c, err := getRandomChar(charset)
		if err != nil {
			return err
		}
		token = append(token, c)
	}

	charset := Uppercase + Lowercase + Number + Symbol
	for i := 4; i < TokenLength; i++ {
		c, err := getRandomChar(charset)
		if err != nil {
			return err
		}
		token = append(token, c)
	}

	if err := shuffle(token); err != nil {
		return err
	}

	v.Code = string(token)

	return nil
}

func getRandomChar(charset string) (byte, error) {
	if len(charset) == 0 {
		return 0, fmt.Errorf("charset must not be empty")
	}
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return 0, fmt.Errorf("GenerateCode error: %s", err)
	}
	return charset[index.Int64()], nil
}

func shuffle(token []byte) error {
	for i := 0; i < len(token)-1; i++ {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(token)-i)))
		if err != nil {
			return fmt.Errorf("GenerateCode error: %s", err)
		}
		token[i], token[int64(i)+j.Int64()] = token[int64(i)+j.Int64()], token[i]
	}
	return nil
}
