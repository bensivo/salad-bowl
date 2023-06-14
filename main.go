package main

import (
	"github.com/bensivo/salad-bowl/game"
)

func main() {
	g := game.NewGame([]string{
		"apple",
		"banana",
		"clementine",
		"durian",
		"fig",
	})

	g.PlayRound()

	// inputChan := make(chan string, 1)
	// defer close(inputChan)

	// Reads input from the terminal, and sends them to the channel
	// go func() {
	// 	for {
	// 		in := bufio.NewReader(os.Stdin)
	// 		result, err := in.ReadString('\n')
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		inputChan <- result
	// 	}
	// }()

}
