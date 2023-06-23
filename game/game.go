package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bensivo/salad-bowl/hub"
)

type Result string

const (
	RESULT_GUESSED Result = "GUESSED"
	RESULT_SKIPPED Result = "SKIPPED"
	RESULT_END     Result = "END"
)

type Game struct {
	Bowl    *Bowl
	Hub     hub.Hub
	Players []*Player
}

func NewGame(hub hub.Hub) *Game {
	bowl := NewBowl()

	return &Game{
		Bowl:    bowl,
		Hub:     hub,
		Players: []*Player{},
	}
}

// Start sets up all listeners and callbacks that the game needs during normal running.
func (g *Game) Start() {
	g.Hub.OnNewConnection(g.AddPlayer)

	g.Hub.OnMessage(func(playerId string, message hub.Message) {
		fmt.Printf("Received message from player %s: %v", playerId, message)

	})
}

// AddPlayer adds a player to the game by id.
// It sends that player their welcome message, and then broadcasts the updated player list.
func (g *Game) AddPlayer(playerId string) {
	fmt.Printf("New player with id %s\n", playerId)

	// Send the player's id to them so they know what it is
	welcomeMsg := hub.Message{
		Event: "notification.player-id",
		Payload: map[string]interface{}{
			"playerId": playerId,
		},
	}
	g.Hub.SendTo(playerId, welcomeMsg)

	player := NewPlayer(playerId)
	g.Players = append(g.Players, player)

	// Broadcast all PlayerIds to all players
	playerList := make([]string, len(g.Players))
	for i := 0; i < len(g.Players); i++ {
		playerList[i] = g.Players[i].Id
	}

	playerListMsg := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": playerList,
		},
	}
	g.Hub.Broadcast(playerListMsg)
}

// PlayRound triggers the main game loop using 2 channels
// The game sends random words to the "words" chan, and receives the result on "results"
// If the word was guessed, it's removed from the bowl. If it is skipped it is not.
//
// If all words have been pulled, "words" will emit empty strings.
// Caller is responsible for sending the 'END' result when the round is over.
func (g *Game) PlayRound() {

	words := make(chan string)
	results := make(chan Result)

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*10))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		Actor(ctx, words, results)
		wg.Done()
	}()

	go func() {
		Dealer(ctx, g.Bowl, words, results)
		wg.Done()
	}()

	wg.Wait()
}

// Dealer pulls words from the bowl and writes them on the words channel
// Then, it reads from the results channel
//
// Exits when the given context is cancelled, or if the bowl is emptied
func Dealer(ctx context.Context, bowl *Bowl, words chan<- string, results <-chan Result) {
	fmt.Println("Starting dealer")
	defer close(words)
	for {
		word := bowl.GetRandomWord()
		if word == "" {
			fmt.Println("Bowl Empty. Games Over.")
			return
		}

		select {
		case words <- word: // If actor is not running, then this blocks
			select {
			case result := <-results:
				fmt.Println("Received res: " + result)
				bowl.RemoveWord(word)
			case <-ctx.Done():
				fmt.Println("Dealer received end signal while listening for response to " + word)
				return
			}

		case <-ctx.Done():
			fmt.Println("Dealer received end signal while sending word " + word)
			return
		}

	}
}

// Actor reads from the words channel, and emits results on the results channel (after waiting)
//
// Exits when the given context is cancelled, or the empty word is received (signalling bowl is empty)
func Actor(ctx context.Context, wordCh chan string, results chan<- Result) {
	fmt.Println("Starting actor")
	defer close(results)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Actor received end signal while waiting for word")
			return

		case word := <-wordCh:
			if word == "" {
				fmt.Println("Received empty word. Game is over")
				return
			}

			fmt.Println("Received Word: " + word)
			time.Sleep(time.Duration(time.Second * 3))

			// NOTE: this actor just always sends RESULT_GUESSED
			// In reality this would come from user input
			select {
			case <-ctx.Done():
				fmt.Println("Actor received end signal while sending result for " + word)
				return
			case results <- RESULT_GUESSED:
				continue
			}

		}
	}
}
