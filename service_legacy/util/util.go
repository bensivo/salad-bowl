package util

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

var generator *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// SeedRand initializes the random number generator with a constant value
// This makes all the random values deterministic, for easier testing
func SeedRand(seed int64) {
	generator = rand.New(rand.NewSource(seed))
}

func RandStringId() string {
	b := make([]rune, 7)
	for i := range b {
		b[i] = letterRunes[generator.Intn(len(letterRunes))]
	}
	b[3] = '-'
	return string(b)
}
