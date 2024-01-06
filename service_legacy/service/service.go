package service

import (
	"fmt"

	"github.com/bensivo/salad-bowl/state"
	"github.com/bensivo/salad-bowl/util"
)

// GameService contains all the functions for managing game instances, and routing events to the appropriate instance
type GameService interface {

	// CRUD Functions, for REST API
	Create() (state.GameState, error)
	Update(id string, gamestate state.GameState) (state.GameState, error)
	Delete(id string) error
	GetAll() ([]state.GameState, error)
	GetOne(id string) (state.GameState, error)

	// Event handlers, for Websockets
	Handle(id string, event state.EventMetadata) error
}

type service struct {
	games map[string]state.GameState
}

var _ GameService = (*service)(nil)

func NewGameService() GameService {
	return &service{
		games: make(map[string]state.GameState),
	}
}

// Create adds a game to the service
func (svc *service) Create() (state.GameState, error) {
	id := util.RandStringId()

	gamestate := state.NewGameState(id)
	svc.games[id] = gamestate

	return gamestate, nil
}

// Update implements GameService.
func (svc *service) Update(id string, gamestate state.GameState) (state.GameState, error) {
	svc.games[id] = gamestate
	return gamestate, nil
}

// Delete removes a game from the servic
func (svc *service) Delete(id string) error {
	delete(svc.games, id)
	return nil
}

// GetAll returns all game
func (svc *service) GetAll() ([]state.GameState, error) {
	games := make([]state.GameState, len(svc.games))

	for _, game := range svc.games {
		games = append(games, game)
	}

	return games, nil
}

// GetOne returns a single game's state by id
func (svc *service) GetOne(id string) (state.GameState, error) {
	// TODO: return error on not found
	return svc.games[id], nil
}

// Handle applies the given event to the relevant game
func (svc *service) Handle(id string, event state.EventMetadata) error {
	game, err := svc.GetOne(id)
	if err != nil {
		fmt.Printf("Failed fetching game %s: %v\n", id, err)
		return err
	}

	game = game.Apply(event)
	svc.Update(id, game)

	return nil
}
