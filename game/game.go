// game defines one instance of the salad bowl being played
package game

import (
	"errors"
	"math/rand"
)

type Result string

const (
	GUESSED Result = "GUESSED"
	SKIPPED Result = "SKIPPED"
	END     Result = "END"
)

type Game struct {
	Words []string
}

func New() *Game {
	return &Game{
		Words: []string{},
	}
}

func (g *Game) AddWord(word string) {
	g.Words = append(g.Words, word)
}

func (g *Game) GetRandomWord() string {
	if len(g.Words) == 0 {
		return ""
	}

	index := rand.Int() % len(g.Words)
	word := g.Words[index]
	return word
}

func (g *Game) RemoveWord(word string) error {
	for i := 0; i < len(g.Words); i++ {
		if g.Words[i] == word {
			// More efficient way to remove an element: swap it with the last element, then decrease size of slice by one
			g.Words[i] = g.Words[len(g.Words)-1]
			g.Words = g.Words[:len(g.Words)-1]

			return nil
		}
	}

	return errors.New("word not found")
}

// PlayRound triggers the main game loop using 2 channels
// The game sends random words to the "words" chan, and receives the result on "results"
// If the word was guessed, it's removed from the bowl. If it is skipped it is not.
//
// If all words have been pulled, "words" will emit empty strings.
// Caller is responsible for sending the 'END' result when the round is over.
func (g *Game) PlayRound(words chan<- string, results <-chan Result) {
	defer close(words)

	for {
		word := g.GetRandomWord()
		words <- word

		res, open := <-results
		if !open {
			return
		}

		if res == GUESSED {
			g.RemoveWord(word)
		}

		if res == END {
			return
		}
	}
}
