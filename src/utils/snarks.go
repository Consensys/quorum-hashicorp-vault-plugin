package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomSeed(n int) []byte {
	b := make([]byte, n)
	newRand().Read(b)
	return b
}

func newRand() *rand.Rand{
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
