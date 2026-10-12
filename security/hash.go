package security

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	charsLen := big.NewInt(int64(len(Chars)))
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", err
		}
		bytes[i] = Chars[num.Int64()]
	}
	return string(bytes), nil
}

func Digest(hashType crypto.Hash, data []byte) ([]byte, error) {
	switch hashType {
	case crypto.SHA256:
		h := crypto.SHA256.New()
		h.Write(data)
		return h.Sum(nil), nil
	case crypto.SHA512:
		h := crypto.SHA512.New()
		h.Write(data)
		return h.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported hash type: %v", hashType)
	}
}

func SHA256String(data string) string {
	h, _ := Digest(crypto.SHA256, []byte(data))
	return hex.EncodeToString(h)
}
