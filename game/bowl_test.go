package game_test

import (
	"fmt"
	"testing"

	"github.com/bensivo/salad-bowl/game"
	"github.com/stretchr/testify/assert"
)

func TestGame_AddAndGetWord(t *testing.T) {
	b := game.NewBowl()
	b.AddWord("apple")

	word := b.GetRandomWord()
	assert.Equal(
		t,
		"apple",
		word,
	)
}

func TestGame_AddAndGetWords(t *testing.T) {
	b := game.NewBowl()
	for i := 0; i < 10; i++ {
		b.AddWord(fmt.Sprintf("%d", i))
	}

	words := []string{}
	for i := 0; i < 10; i++ {
		word := b.GetRandomWord()
		words = append(words, word)
		b.RemoveWord(word)
	}

	for i := 0; i < 10; i++ {
		assert.Contains(t, words, fmt.Sprintf("%d", i))
	}
}

func TestGame_RemoveWord(t *testing.T) {
	b := game.NewBowl()

	b.AddWord("hello")
	b.AddWord("world")

	assert.Equal(t, []string{"hello", "world"}, b.Words)
	b.RemoveWord("world")

	assert.Equal(t, []string{"hello"}, b.Words)
}
