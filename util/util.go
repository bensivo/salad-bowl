package util

import (
	"math/rand"
)

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var generator *rand.Rand

func Init() {
	generator = rand.New(rand.NewSource(0)) // Seed with a 0 for deterministic testing
	// generator = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandStringId() string {
	b := make([]rune, 7)
	for i := range b {
		b[i] = letterRunes[generator.Intn(len(letterRunes))]
	}
	b[3] = '-'
	return string(b)
}
