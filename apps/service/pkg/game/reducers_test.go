package game_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/game/db"
	"github.com/bensivo/salad-bowl/service/pkg/util"
)

func TestHandlePlayerJoinedAddsPlayer(t *testing.T) {
	// Given a game exists
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// When that game receives a playerjoined event
	g.HandlePlayerJoined(game.PlayerJoinedEventPayload{
		PlayerID:   "11111",
		PlayerName: "Alice",
	})

	// Then the player is added to the game
	found := false
	for _, player := range g.Players {
		if player.PlayerID == "11111" {
			found = true
			break
		}
	}
	if found == false {
		t.Errorf("player 11111 not found in game state")
	}
}

func TestHandlePlayerJoinedSkipsDuplicates(t *testing.T) {
	// Given a game exists, and a player was added already
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// Given a player was added already
	g.HandlePlayerJoined(game.PlayerJoinedEventPayload{
		PlayerID:   "11111",
		PlayerName: "Alice",
	})

	// When we send another playerjoined event
	g.HandlePlayerJoined(game.PlayerJoinedEventPayload{
		PlayerID:   "11111",
		PlayerName: "Alice",
	})

	// Then the player is not added twice
	count := 0
	for _, player := range g.Players {
		if player.PlayerID == "11111" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("player 11111 was added twice to game state")
	}
}

func TestHandlePlayerLeftRemovesPlayer(t *testing.T) {
	// Given a game exists
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// Given a player is in the game
	g.HandlePlayerJoined(game.PlayerJoinedEventPayload{
		PlayerID:   "11111",
		PlayerName: "Alice",
	})

	// When the game recieves a player-left event
	g.HandlePlayerLeft(game.PlayerLeftEventPayload{
		PlayerID: "11111",
	})

	// Then the player is removed from the game
	found := slices.ContainsFunc(g.Players, func(p game.Player) bool {
		return p.PlayerID == "11111"
	})
	if found == true {
		t.Errorf("player 11111 is still in game state")
	}
}

func TestHandlePlayerLeftRemovesPlayerFromTeam(t *testing.T) {
	// Given a game exists
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// Given a player is in the game
	g.HandlePlayerJoined(game.PlayerJoinedEventPayload{
		PlayerID:   "11111",
		PlayerName: "Alice",
	})

	// Given the player is in a team
	g.HandleTeamJoined(game.TeamJoinedEventPayload{
		PlayerID: "11111",
		TeamName: "Red",
	})

	// When the game recieves a player-left event
	g.HandlePlayerLeft(game.PlayerLeftEventPayload{
		PlayerID: "11111",
	})

	// Then the player is removed from the team
	g, _ = gameSvc.GetOne(g.ID)
	if len(g.Teams[0].PlayerIDs) != 0 {
		t.Errorf("team not emptied after player left")
	}
}

func TestHadnleTeamJoinedAddsPlayerToTeam(t *testing.T) {
	// Given a game exists
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// Given a player is in the game
	g.HandlePlayerJoined(game.PlayerJoinedEventPayload{
		PlayerID:   "11111",
		PlayerName: "Alice",
	})

	// Given the player is in a team
	g.HandleTeamJoined(game.TeamJoinedEventPayload{
		PlayerID: "11111",
		TeamName: "Red",
	})

	// Then the player is added to that team's roster
	if g.Teams[0].PlayerIDs[0] != "11111" {
		t.Errorf("player 11111 not found in Red team")
	}
}

func TestHandleTeamJoinedPlayerNotFound(t *testing.T) {
	// Given a game exists, without any players
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// When we send team-joined
	err := g.HandleTeamJoined(game.TeamJoinedEventPayload{
		PlayerID: "11111",
		TeamName: "Red",
	})

	if !errors.Is(err, game.ErrPlayerNotFound) {
		t.Errorf("expected error not returned")
	}
}

func TestHandleTeamJoinedRemovesPlayerFromOldTeam(t *testing.T) {
	// Given a game exists
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// Given a player is in the game
	g.HandlePlayerJoined(game.PlayerJoinedEventPayload{
		PlayerID:   "11111",
		PlayerName: "Alice",
	})

	// Given the player is the Red team
	g.HandleTeamJoined(game.TeamJoinedEventPayload{
		PlayerID: "11111",
		TeamName: "Red",
	})

	// Wehn the player joins the Blue team
	g.HandleTeamJoined(game.TeamJoinedEventPayload{
		PlayerID: "11111",
		TeamName: "Blue",
	})

	// Then the game shows the player in the blue team, and not the red team
	if slices.Contains(g.Teams[0].PlayerIDs, "11111") {
		t.Errorf("player 11111 is not supposed to be in the red team")
	}
	if !slices.Contains(g.Teams[1].PlayerIDs, "11111") {
		t.Errorf("player 11111 not found in the blue team")
	}
}
