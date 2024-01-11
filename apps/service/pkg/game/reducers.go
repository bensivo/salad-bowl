// reducers.go is where the bulk of salad bowl's business logic lives.
// This file defines handler functions for each GameEvent, and how game state updates based on each event's payload.
package game

import (
	"slices"

	"github.com/bensivo/salad-bowl/service/pkg/log"
)

func (g *Game) HandlePlayerJoined(payload PlayerJoinedEventPayload) {
	log.Infof("Received player joined event %s(%s)\n", payload.PlayerName, payload.PlayerID)

	// Check for duplicates
	for _, player := range g.Players {
		if player.PlayerID == payload.PlayerID {
			log.Infof("Player %s has already joined the game\n", payload.PlayerID)
			return
		}
	}

	// Add player to game state
	g.Players = append(g.Players, Player{
		PlayerID:   payload.PlayerID,
		PlayerName: payload.PlayerName,
	})
}

func (g *Game) HandlePlayerLeft(payload PlayerLeftEventPayload) {
	log.Infof("Received player left event %s\n", payload.PlayerID)

	// Remove player from game.Players
	g.Players = slices.DeleteFunc(g.Players, func(p Player) bool {
		return p.PlayerID == payload.PlayerID
	})

	// Remove player from all teams
	for i := range g.Teams {
		g.Teams[i].PlayerIDs = slices.DeleteFunc(g.Teams[i].PlayerIDs, func(playerId string) bool {
			return playerId == payload.PlayerID
		})
	}
}

func (g *Game) HandleTeamJoined(payload TeamJoinedEventPayload) error {
	log.Infof("Received team joined event Player:%s Team:%s\n", payload.PlayerID, payload.TeamName)

	// Make sure player exists
	exists := slices.ContainsFunc(g.Players, func(p Player) bool {
		return p.PlayerID == payload.PlayerID
	})
	if !exists {
		log.Infof("Player %s does not exist\n", payload.PlayerID)
		return ErrPlayerNotFound
	}

	for i, team := range g.Teams {
		// Remove player from all teams, preventing a player from being in 2 at once.
		g.Teams[i].PlayerIDs = slices.DeleteFunc(g.Teams[i].PlayerIDs, func(playerId string) bool {
			return playerId == payload.PlayerID
		})

		// Add player to the team which matches the event payload
		if team.TeamName == payload.TeamName {
			log.Infof("Adding player %s to team %d\n", payload.PlayerID, i)
			g.Teams[i].PlayerIDs = append(team.PlayerIDs, payload.PlayerID)
		}
	}

	return nil
}
