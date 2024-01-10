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

type MockListener struct {
	calls []game.Game
}

func (ml *MockListener) OnChange(g game.Game) {
	ml.calls = append(ml.calls, g)
}

func TestRegisterListenerGetsCalledWithCurrentState(t *testing.T) {
	// Given a game exists
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// Given some events have already been handled, causing game state to change
	gameSvc.HandleEvent(g.ID, game.GameEvent{
		Name:      game.PlayerJoined,
		Timestamp: util.NowIso8601(),
		Payload: game.PlayerJoinedEventPayload{
			PlayerID:   "11111",
			PlayerName: "alice",
		},
	})

	// When a listener is registered
	listener := &MockListener{
		calls: make([]game.Game, 0),
	}
	gameSvc.RegisterListener(g.ID, listener)

	// Then the listener is called immediately with the current game state
	if len(listener.calls) != 1 {
		t.Errorf("expected listener to be called")
	}
}
func TestRegisterListenerReceivesChanges(t *testing.T) {
	// Given a game exists
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)
	g, _ := gameSvc.Create()

	// Given a listener registered for that game
	listener := &MockListener{
		calls: make([]game.Game, 0),
	}
	gameSvc.RegisterListener(g.ID, listener)

	// When an event causes a change in game state
	gameSvc.HandleEvent(g.ID, game.GameEvent{
		Name:      game.PlayerJoined,
		Timestamp: util.NowIso8601(),
		Payload: game.PlayerJoinedEventPayload{
			PlayerID:   "11111",
			PlayerName: "alice",
		},
	})

	// Then the listener is called
	if len(listener.calls) != 2 { // calls = 2 because there is one call when the listener first registers, and another after the event
		t.Errorf("expected listener to be called")
	}

	// Then the call has the change
	if listener.calls[1].Players[0].PlayerID != "11111" {
		t.Errorf("expected player to be added to the game")
	}
}

// TODO:
//   - when a game is deleted, call the close method on any listeners
//   - when deregister is called, stop calling that listener
