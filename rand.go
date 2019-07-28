package main

import (
	cryptorand "crypto/rand"
	"math/big"
	pseudorand "math/rand"
	"time"
)

const MaxInt64 = 0x7fffffffffffffff

func RandomSeed() int64 {
	if rand, err := cryptorand.Int(cryptorand.Reader, big.NewInt(MaxInt64)); err == nil {
		return rand.Int64()
	}
	pseudorand.Seed(pseudorand.Int63() ^ time.Now().UnixNano())
	return pseudorand.Int63()
}

func Randomize(seed int64) {
	if seed == 0 {
		seed = RandomSeed()
	}
	Logf("Global random seed: %d", seed)
	pseudorand.Seed(seed)
}

func NewRand(seed int64) *pseudorand.Rand {
	return pseudorand.New(pseudorand.NewSource(seed))
}
