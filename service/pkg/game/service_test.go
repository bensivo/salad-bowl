package game_test

import (
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
}
