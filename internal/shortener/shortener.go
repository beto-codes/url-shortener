package shortener

import (
	"crypto/md5"
	"encoding/base64"
	"math/rand"
	"time"
)

const (
	shortCodeLength = 6
)

type Generator struct {
	rng *rand.Rand
}

func NewGenerator() *Generator {
	return &Generator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Generator) Generate(longURL string) string {
	hash := md5.Sum([]byte(longURL))
	encoded := base64.URLEncoding.EncodeToString(hash[:])

	if len(encoded) > shortCodeLength {
		encoded = encoded[:shortCodeLength]
	}

	return encoded
}

func (g *Generator) GenerateRandom() string {
	const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, shortCodeLength)
	for i := range b {
		b[i] = charSet[g.rng.Intn(len(charSet))]
	}

	return string(b)
}
