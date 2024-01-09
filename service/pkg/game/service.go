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
		Phase:            "lobby",
		SubmittedWords:   []SubmittedWord{},
		RemainingWords:   []string{},
		RemainingPlayers: []string{},
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

		err = svc.db.Save(game)
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

		// Remove player from game.Players
		game.Players = slices.DeleteFunc(game.Players, func(p Player) bool {
			return p.PlayerID == payload.PlayerID
		})

		// Remove player from all teams
		for i := range game.Teams {
			game.Teams[i].PlayerIDs = slices.DeleteFunc(game.Teams[i].PlayerIDs, func(playerId string) bool {
				return playerId == payload.PlayerID
			})
		}

		err = svc.db.Save(game)
		if err != nil {
			return fmt.Errorf("failed updating game state %v", err)
		}
	case TeamJoined:
		var payload TeamJoinedEventPayload
		err := ParseGameEventPayload(event.Payload, &payload)
		if err != nil {
			return fmt.Errorf("failed parsing event payload: %v", err)
		}

		log.Infof("Received team joined event Player:%s Team:%s\n", payload.PlayerID, payload.TeamName)

		// Make sure player exists
		exists := slices.ContainsFunc(game.Players, func(p Player) bool {
			return p.PlayerID == payload.PlayerID
		})
		if !exists {
			log.Infof("Player %s does not exist\n", payload.PlayerID)
			return ErrPlayerNotFound
		}

		for i, team := range game.Teams {

			// Remove player from all teams, preventing a player from being in 2 at once.
			game.Teams[i].PlayerIDs = slices.DeleteFunc(game.Teams[i].PlayerIDs, func(playerId string) bool {
				return playerId == payload.PlayerID
			})

			if team.TeamName == payload.TeamName {
				log.Infof("Adding player %s to team %d\n", payload.PlayerID, i)
				game.Teams[i].PlayerIDs = append(team.PlayerIDs, payload.PlayerID)
			}
		}

		// TODO: handle player which is already in a team, leaving it and joining another team

		err = svc.db.Save(game)
		if err != nil {
			return fmt.Errorf("failed updating game state %v", err)
		}
	}

	return nil
}
