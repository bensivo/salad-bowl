package game_test

import (
	"errors"
	"testing"

	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/game/db"
	"github.com/bensivo/salad-bowl/service/pkg/util"
)

func TestCreateGeneratesRandomId(t *testing.T) {
	util.SeedRand(0)
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)

	g, _ := gameSvc.Create()

	if g.ID != "CUB-HIZ" {
		t.Errorf("Expected ID to be CUB-HIZ from seeded random generator")
	}
}

func TestHandleEvent(t *testing.T) {

	t.Run("PlayerJoinedEvent - adds player", func(t *testing.T) {
		// Given a game exists
		util.SeedRand(0)
		gameDb := db.NewInMemoryGameDb()
		gameSvc := game.NewGameService(gameDb)
		g, _ := gameSvc.Create()

		// When that game receives a playerjoined event
		e := game.GameEvent{
			Name:      game.PlayerJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.PlayerJoinedEventPayload{
				PlayerID:   "11111",
				PlayerName: "Alice",
			},
		}
		gameSvc.HandleEvent(g.ID, e)

		// Then the player is added to the game
		g, _ = gameSvc.GetOne(g.ID)
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
	})

	t.Run("PlayerJoinedEvent - duplicate player id", func(t *testing.T) {
		// Given a game exists, and a player was added already
		util.SeedRand(0)
		gameDb := db.NewInMemoryGameDb()
		gameSvc := game.NewGameService(gameDb)
		g, _ := gameSvc.Create()
		e := game.GameEvent{
			Name:      game.PlayerJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.PlayerJoinedEventPayload{
				PlayerID:   "11111",
				PlayerName: "Alice",
			},
		}
		gameSvc.HandleEvent(g.ID, e)

		// When we send another playerjoined event
		gameSvc.HandleEvent(g.ID, e)

		// Then the player is not added twice
		g, _ = gameSvc.GetOne(g.ID)
		count := 0
		for _, player := range g.Players {
			if player.PlayerID == "11111" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("player 11111 was added twice to game state")
		}
	})

	t.Run("PlayerLeftEvent - removes player", func(t *testing.T) {
		// Given a game exists, and a player is in that game
		util.SeedRand(0)
		gameDb := db.NewInMemoryGameDb()
		gameSvc := game.NewGameService(gameDb)
		g, _ := gameSvc.Create()
		gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.PlayerJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.PlayerJoinedEventPayload{
				PlayerID:   "11111",
				PlayerName: "Alice",
			},
		})

		// When the game recieves a player-left event
		gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.PlayerLeft,
			Timestamp: util.NowIso8601(),
			Payload: game.PlayerLeftEventPayload{
				PlayerID: "11111",
			},
		})

		// Then the player is removed from the game
		g, _ = gameSvc.GetOne(g.ID)
		found := false
		for _, player := range g.Players {
			if player.PlayerID == "11111" {
				found = true
				break
			}
		}
		if found == true {
			t.Errorf("player 11111 is still in game state")
		}
	})

	t.Run("PlayerLeftEvent - removes player from team", func(t *testing.T) {
		// Given a game exists, and a player is in that game, and the player is in a team
		util.SeedRand(0)
		gameDb := db.NewInMemoryGameDb()
		gameSvc := game.NewGameService(gameDb)
		g, _ := gameSvc.Create()

		gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.PlayerJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.PlayerJoinedEventPayload{
				PlayerID:   "11111",
				PlayerName: "Alice",
			},
		})
		gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.TeamJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.TeamJoinedEventPayload{
				PlayerID: "11111",
				TeamName: "Red",
			},
		})

		// When the game recieves a player-left event
		gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.PlayerLeft,
			Timestamp: util.NowIso8601(),
			Payload: game.PlayerLeftEventPayload{
				PlayerID: "11111",
			},
		})

		// Then the player is removed from the team
		g, _ = gameSvc.GetOne(g.ID)
		if len(g.Teams[0].PlayerIDs) != 0 {
			t.Errorf("team not emptied after player left")
		}
	})

	t.Run("TeamJoinedEvent - adds player to team", func(t *testing.T) {
		// Given a game exists, with a player
		util.SeedRand(0)
		gameDb := db.NewInMemoryGameDb()
		gameSvc := game.NewGameService(gameDb)
		g, _ := gameSvc.Create()
		gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.PlayerJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.PlayerJoinedEventPayload{
				PlayerID:   "11111",
				PlayerName: "Alice",
			},
		})

		// When we send team-joined
		gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.TeamJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.TeamJoinedEventPayload{
				PlayerID: "11111",
				TeamName: "Red",
			},
		})

		// Then the player is added to that team's roster
		g, _ = gameSvc.GetOne(g.ID)
		if g.Teams[0].PlayerIDs[0] != "11111" {
			t.Errorf("player 11111 not found in Red team")
		}
	})

	t.Run("TeamJoinedEvent - nonexistent player", func(t *testing.T) {
		// Given a game exists, without any players
		util.SeedRand(0)
		gameDb := db.NewInMemoryGameDb()
		gameSvc := game.NewGameService(gameDb)
		g, _ := gameSvc.Create()

		// When we send team-joined
		err := gameSvc.HandleEvent(g.ID, game.GameEvent{
			Name:      game.TeamJoined,
			Timestamp: util.NowIso8601(),
			Payload: game.TeamJoinedEventPayload{
				PlayerID: "11111",
				TeamName: "Red",
			},
		})

		if !errors.Is(err, game.ErrPlayerNotFound) {
			t.Errorf("expected error not returned")
		}
	})

}
