package game_test

import (
	"fmt"
	"testing"

	"github.com/bensivo/salad-bowl/game"
	"github.com/stretchr/testify/assert"
)

func TestGame_AddAndGetWord(t *testing.T) {
	g := game.New()
	g.AddWord("apple")

	word := g.GetRandomWord()
	assert.Equal(
		t,
		"apple",
		word,
	)
}

func TestGame_AddAndGetWords(t *testing.T) {
	g := game.New()
	for i := 0; i < 10; i++ {
		g.AddWord(fmt.Sprintf("%d", i))
	}

	words := []string{}
	for i := 0; i < 10; i++ {
		word := g.GetRandomWord()
		words = append(words, word)
		g.RemoveWord(word)
	}

	for i := 0; i < 10; i++ {
		assert.Contains(t, words, fmt.Sprintf("%d", i))
	}
}

func TestGame_RemoveWord(t *testing.T) {
	g := game.New()

	g.AddWord("hello")
	g.AddWord("world")

	assert.Equal(t, []string{"hello", "world"}, g.Words)
	g.RemoveWord("world")

	assert.Equal(t, []string{"hello"}, g.Words)
}

func TestGame_Start_RemovesWordsOnGuessed(t *testing.T) {
	g := game.New()

	g.AddWord("hello")

	words := make(chan string)
	results := make(chan game.Result)
	go func() {
		defer close(results)
		word := <-words
		assert.Equal(t, word, "hello")

		results <- game.GUESSED

		word = <-words
		assert.Equal(t, word, "")
		results <- game.END
	}()

	assert.Equal(t, g.Words, []string{"hello"})
	g.PlayRound(words, results)
	assert.Equal(t, g.Words, []string{})
}

func TestGame_Start_KeepsWordsOnSkipped(t *testing.T) {
	g := game.New()

	g.AddWord("hello")

	words := make(chan string)
	results := make(chan game.Result)

	go func() {
		defer close(results)
		word := <-words
		assert.Equal(t, word, "hello")

		results <- game.SKIPPED

		word = <-words
		assert.Equal(t, word, "hello")
		results <- game.END
	}()

	assert.Equal(t, g.Words, []string{"hello"})
	g.PlayRound(words, results)
	assert.Equal(t, g.Words, []string{"hello"})
}
