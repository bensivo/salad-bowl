package game

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/bensivo/salad-bowl/service/pkg/log"
	"github.com/bensivo/salad-bowl/service/pkg/util"
)

var ErrPlayerNotFound = errors.New("player not found")

type service struct {
	db        GameDb
	listeners map[string][]GameListener
}

var _ GameService = (*service)(nil)

func NewGameService(db GameDb) GameService {
	return &service{
		db:        db,
		listeners: map[string][]GameListener{},
	}
}

// Create implements GameService.
func (svc *service) Create() (*Game, error) {
	id := util.RandStringId()
	g := &Game{
		ID:        id,
		CreatedAt: time.Now(),
		Players:   []Player{},
		Teams: []Team{
			{
				TeamName:  "Red",
				Score:     0,
				PlayerIDs: []string{},
			},
			{
				TeamName:  "Blue",
				Score:     0,
				PlayerIDs: []string{},
			},
			{
				TeamName:  "Spectator",
				Score:     0,
				PlayerIDs: []string{},
			},
		},
		Phase: "lobby",
		// SubmittedWords:   []SubmittedWord{},
		// RemainingWords:   []string{},
		// RemainingPlayers: []string{},
	}
	err := svc.db.Save(g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// GetAll implements GameService.
func (svc *service) GetAll() ([]*Game, error) {
	return svc.db.GetAll()
}

// GetOne implements GameService.
func (svc *service) GetOne(ID string) (*Game, error) {
	return svc.db.GetOne(ID)
}

// Delete implements GameService.
func (svc *service) Delete(ID string) error {

	// TODO: clean up the listeners map, remove any listeners for this game before deleting it

	return svc.db.Delete(ID)
}

// HandleEvent implements GameService.
func (svc *service) HandleEvent(gameID string, event GameEvent) error {
	// Find the game related to this event
	game, err := svc.GetOne(gameID)
	if err != nil {
		return fmt.Errorf("failed fetching game %s for event: %v", gameID, err)
	}

	// Based on the event name, call the appropriate handler function from reducers.go
	switch event.Name {
	case PlayerJoined:
		var payload PlayerJoinedEventPayload
		err := ParseGameEventPayload(event.Payload, &payload)
		if err != nil {
			return fmt.Errorf("failed parsing event payload: %v", err)
		}

		game.HandlePlayerJoined(payload)

	case PlayerLeft:
		var payload PlayerLeftEventPayload
		err := ParseGameEventPayload(event.Payload, &payload)
		if err != nil {
			return fmt.Errorf("failed parsing event payload: %v", err)
		}

		game.HandlePlayerLeft(payload)

	case TeamJoined:
		var payload TeamJoinedEventPayload
		err := ParseGameEventPayload(event.Payload, &payload)
		if err != nil {
			return fmt.Errorf("failed parsing event payload: %v", err)
		}

		game.HandleTeamJoined(payload)
	}

	// Save the new game state to the database
	err = svc.db.Save(game)
	if err != nil {
		return fmt.Errorf("failed updating game state %v", err)
	}

	listeners, found := svc.listeners[game.ID]
	if found {
		for _, listener := range listeners {
			listener.OnChange(*game)
		}
	}

	return nil
}

// RegisterListener implements GameService.
func (svc *service) RegisterListener(ID string, listener GameListener) error {
	// Find the game related to this event
	game, err := svc.GetOne(ID)
	if err != nil {
		log.Infof("Error getting game: %v\n", err)
		return fmt.Errorf("failed fetching game %s for event: %v", ID, err)
	}

	// Add the listener to the listeners list
	_, found := svc.listeners[ID]
	if !found {
		svc.listeners[ID] = []GameListener{}
	}
	svc.listeners[ID] = append(svc.listeners[ID], listener)

	// Call the listener with the current game state
	listener.OnChange(*game)
	return nil
}

// DeregisterListener implements GameService.
func (svc *service) DeregisterListener(ID string, listener GameListener) error {
	_, found := svc.listeners[ID]
	if !found {
		return nil // TODO: should we return nil? Is it an error to deregister from a non-existent game?
	}

	svc.listeners[ID] = slices.DeleteFunc(svc.listeners[ID], func(l GameListener) bool {
		return l == listener
	})

	return nil
}
