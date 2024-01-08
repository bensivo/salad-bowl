package game

import (
	"fmt"
	"slices"

	"github.com/bensivo/salad-bowl/service/pkg/log"
	"github.com/bensivo/salad-bowl/service/pkg/util"
)

type service struct {
	db GameDb
}

var _ GameService = (*service)(nil)

func NewGameService(db GameDb) GameService {
	return &service{
		db: db,
	}
}

// Create implements GameService.
func (svc *service) Create() (*Game, error) {
	id := util.RandStringId()
	g, err := svc.db.Create(id)
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
	return svc.db.Delete(ID)
}

// HandleEvent implements GameService.
func (svc *service) HandleEvent(gameID string, event GameEvent) error {
	// Find the game related to this event
	game, err := svc.GetOne(gameID)
	if err != nil {
		return fmt.Errorf("failed fetching game %s for event: %v", gameID, err)
	}

	switch event.Name {
	case PlayerJoined:
		var payload PlayerJoinedEventPayload
		err := ParseGameEventPayload(event.Payload, &payload)
		if err != nil {
			return fmt.Errorf("failed parsing event payload: %v", err)
		}

		log.Infof("Received player joined event %s(%s)\n", payload.PlayerName, payload.PlayerID)

		// Check for duplicates
		for _, player := range game.Players {
			if player.PlayerID == payload.PlayerID {
				log.Infof("Player %s has already joined the game\n", payload.PlayerID)
				return nil
			}
		}

		// Add player to game state
		game.Players = append(game.Players, Player{
			PlayerID:   payload.PlayerID,
			PlayerName: payload.PlayerName,
		})

		err = svc.db.Update(game.ID, game)
		if err != nil {
			return fmt.Errorf("failed updating game state %v", err)
		}
	case PlayerLeft:
		var payload PlayerLeftEventPayload
		err := ParseGameEventPayload(event.Payload, &payload)
		if err != nil {
			return fmt.Errorf("failed parsing event payload: %v", err)
		}

		log.Infof("Received player left event %s\n", payload.PlayerID)

		game.Players = slices.DeleteFunc(game.Players, func(p Player) bool {
			return p.PlayerID == payload.PlayerID
		})

		err = svc.db.Update(game.ID, game)
		if err != nil {
			return fmt.Errorf("failed updating game state %v", err)
		}
	}

	return nil
}
