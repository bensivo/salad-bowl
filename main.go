package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bensivo/salad-bowl/game"
)

func main() {
	g := game.New()

	g.AddWord("apple")
	g.AddWord("banana")
	g.AddWord("clementine")
	g.AddWord("durian")
	g.AddWord("fig")

	inputChan := make(chan string, 1)
	defer close(inputChan)

	// Reads input from the terminal, and sends them to the channel
	go func() {
		for {
			in := bufio.NewReader(os.Stdin)
			result, err := in.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			inputChan <- result
		}
	}()

	playRound := func() {
		words := make(chan string)
		results := make(chan game.Result)
		timer := time.After(time.Second * 10)

		go func() {
			for {
				word := <-words
				if word == "" {
					results <- game.END
				}

				fmt.Printf("Your word is: %s\n", word)

				select {
				case input := <-inputChan:
					if strings.ToLower(input)[0] == 'y' {
						fmt.Printf("Confirmed word: '%s'\n", word)
						results <- game.GUESSED
					} else if strings.ToLower(input)[0] == 'n' {
						fmt.Printf("Putting word '%s' back in the bowl\n", word)
						results <- game.SKIPPED
					} else {
						fmt.Println("Please enter either 'y' or 'n'")
					}
				case <-timer:
					fmt.Printf("Times up! There are %d words left\n", len(g.Words))
					results <- game.END
					return
				}
			}
		}()

		g.PlayRound(words, results)
	}

	for i := 0; i < 3; i++ {
		fmt.Printf("--------------\n")
		fmt.Printf("ROUND %d\n", i)
		fmt.Printf("--------------\n")

		// TODO: Naming - rounds only end when the word bank is empty
		// The timer is just for each individual

		// TODO: add a concept of points, how many did you remove? And what team were you on?
		playRound()
	}
}
