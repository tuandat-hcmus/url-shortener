package utils

import (
	"crypto/rand"
	"math/big"
)

//go:generate mockery --name CodeGen --filename codegen.go --outpkg utils_mock
type CodeGen interface {
	GenerateCode(length int) (string, error)
}

type codeGen struct{}

func NewCodeGen() CodeGen {
	return &codeGen{}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateCode generates a random code of the specified length
func (c *codeGen) GenerateCode(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result), nil
}
