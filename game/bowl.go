package game

import (
	"errors"
	"math/rand"
)

type Bowl struct {
	Words []string // All words available in the game
}

func NewBowl() *Bowl {
	return &Bowl{
		Words: []string{},
	}
}

func (g *Bowl) AddWord(word string) {
	g.Words = append(g.Words, word)
}

func (g *Bowl) GetRandomWord() string {
	if len(g.Words) == 0 {
		return ""
	}

	index := rand.Int() % len(g.Words)
	word := g.Words[index]
	return word
}

func (g *Bowl) RemoveWord(word string) error {
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
