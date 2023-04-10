package util

import (
	cryptoRand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
)

const randomChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	seed, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
}

// RandomString ランダムな文字列を生成
func RandomString(n int) string {
	b := make([]byte, n)
	l := len(randomChars)
	for i := range b {
		b[i] = randomChars[rand.Intn(l)]
	}
	return string(b)
}
