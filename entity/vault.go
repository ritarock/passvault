package entity

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase = "abcdefghijklmnopqrstuvwxyz"
	Number    = "0123456789"
	Symbol    = "!#$%&()-@_<>"
)

const TokenLength = 12

type Vault struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Code  string `json:"code"`
	Iv    string `json:"iv"`
}

func (v *Vault) GenerateCode() error {
	var token []byte
	for _, charset := range []string{Uppercase, Lowercase, Number, Symbol} {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return fmt.Errorf("generateCode error: %s", err)
		}
		token = append(token, charset[randomIndex.Int64()])
	}

	charset := Uppercase + Lowercase + Number + Symbol
	for i := 4; i < TokenLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return fmt.Errorf("generateCode error: %s", err)
		}
		token = append(token, charset[randomIndex.Int64()])
	}

	for i := 0; i < len(token)-1; i++ {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(token)-i)))
		if err != nil {
			return fmt.Errorf("generateCode error: %s", err)
		}
		token[i], token[int64(i)+j.Int64()] = token[int64(i)+j.Int64()], token[i]
	}

	v.Code = string(token)

	return nil
}
