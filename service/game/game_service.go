package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/bensivo/salad-bowl/hub"
)

type GameService struct {
	games map[string]*Game
}

func NewGameService() *GameService {
	return &GameService{
		games: make(map[string]*Game),
	}
}

func (svc *GameService) Create() (string, error) {
	if len(svc.games) >= 10 {
		fmt.Printf("Could not create game. Instance full.\n")
		return "", errors.New("instance full. There are already 10 games in the server")
	}

	newHub := hub.NewHub()

	newGame := NewGame(newHub)
	newGame.Start()

	gameId := newGame.ID
	svc.games[gameId] = newGame

	fmt.Printf("Created new game with id: %s\n", gameId)

	return gameId, nil
}

func (svc *GameService) Delete(id string) {
	fmt.Printf("Deleting game %s\n", id)
	delete(svc.games, id)
}

func (svc *GameService) GetOne(id string) (*Game, error) {
	l, ok := svc.games[id]
	if !ok {
		return nil, fmt.Errorf("game %s not found", id)
	}

	return l, nil
}

func (svc *GameService) Get() map[string]*Game {
	return svc.games
}

// Cleanup deletes all games that have no players, and are more than 30 seconds old
func (svc *GameService) Cleanup() {
	now := time.Now()
games:
	for id, game := range svc.games {
		// Don't remove games with at least 1 online player
		for _, player := range game.Players.Get() {
			if player.Status == "online" {
				continue games
			}
		}

		// Don't remove brand new games. Here, "brand new" was arbitrarily set to 30 seconds. May realize this needs to be longer in the future.
		if now.Sub(game.CreatedAt) < time.Duration(30*time.Second) {
			continue games
		}

		fmt.Printf("Deleting empty game %s, created at %s\n", id, game.CreatedAt)
		delete(svc.games, id)
	}
}

// StartCleanup schedules the cleanup job to run every 10 seconds
func (svc *GameService) StartCleanup() {
	go func() {
		for {
			time.Sleep(time.Duration(10 * time.Second))
			svc.Cleanup()
		}
	}()
}
