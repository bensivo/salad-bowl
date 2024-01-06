package state_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/state"
	"github.com/stretchr/testify/assert"
)

func Test_NewGameState_SetsId(t *testing.T) {
	g := state.NewGameState("game-id")

	assert.Equal(t, g.ID, "game-id")
}

// func Test_Apply_PlayerJoined_AddsPlayer(t *testing.T) {
// 	// Given a new game
// 	g := state.NewGameState("game-id")

// 	// When a PlayerJoinedEvent is received
// 	g = g.Apply(state.PlayerJoinedEvent{
// 		PlayerID:   "player-1",
// 		PlayerName: "alice",
// 	})

// 	// Then a new player is added to the game
// 	assert.Equal(t, g.Players, []state.Player{
// 		{
// 			PlayerID:   "player-1",
// 			PlayerName: "alice",
// 			TeamName:   "",
// 		},
// 	})
// }

// func Test_Apply_PlayerJoined_IgnoresDuplicates(t *testing.T) {
// 	// Given a new game, with a player
// 	g := state.NewGameState("game-id")
// 	g = g.Apply(state.PlayerJoinedEvent{
// 		PlayerID:   "player-1",
// 		PlayerName: "alice",
// 	})

// 	// When a duplicate PlayerJoinedEvent is received
// 	g = g.Apply(state.PlayerJoinedEvent{
// 		PlayerID:   "player-1",
// 		PlayerName: "alice",
// 	})

// 	// Then the new player is NOT added to the game
// 	assert.Equal(t, g.Players, []state.Player{
// 		{
// 			PlayerID:   "player-1",
// 			PlayerName: "alice",
// 			TeamName:   "",
// 		},
// 	})
// }

// func Test_Apply_PlayerLeft_RemovesPlayer(t *testing.T) {
// 	// Given a new Game, with 2 players
// 	g := state.NewGameState("game-id")
// 	g = g.Apply(state.PlayerJoinedEvent{
// 		PlayerID:   "player-1",
// 		PlayerName: "alice",
// 	})
// 	g = g.Apply(state.PlayerJoinedEvent{
// 		PlayerID:   "player-2",
// 		PlayerName: "bob",
// 	})

// 	// When the player removed event is received
// 	g = g.Apply((state.PlayerLeftEvent{
// 		PlayerID: "player-1",
// 	}))

// 	// Then the player is removed
// 	assert.Equal(t, g.Players, []state.Player{
// 		{
// 			PlayerID:   "player-2",
// 			PlayerName: "bob",
// 			TeamName:   "",
// 		},
// 	})

// 	// When the last player is removed
// 	g = g.Apply((state.PlayerLeftEvent{
// 		PlayerID: "player-2",
// 	}))

// 	// Then the player list is empty
// 	assert.Equal(t, g.Players, []state.Player{})
// }

// func Test_Apply_TeamJoined_UpdatedPlayerTeam(t *testing.T) {
// 	// Given a new Game, with a player
// 	g := state.NewGameState("game-id")
// 	g = g.Apply(state.PlayerJoinedEvent{
// 		PlayerID:   "player-1",
// 		PlayerName: "alice",
// 	})

// 	// When the player sends TeamJoined
// 	g = g.Apply((state.TeamJoinedEvent{
// 		PlayerID: "player-1",
// 		TeamName: "team-1",
// 	}))

// 	// Then the player's team is updated
// 	assert.Equal(t, g.Players, []state.Player{
// 		{
// 			PlayerID:   "player-1",
// 			PlayerName: "alice",
// 			TeamName:   "team-1",
// 		},
// 	})
// }

// func Test_Apply_WordBankStarted_UpdatesPhase(t *testing.T) {
// 	// Given a new Game, with a player
// 	g := state.NewGameState("game-id")

// 	// When we apply the WordBankStarted event
// 	g = g.Apply((state.WordBankStartedEvent{}))

// 	// Then the phase is updated
// 	assert.Equal(t, g.Phase, "word-bank")
// }

// func Test_Apply_WordAdded_UpdatesSubmittedWords(t *testing.T) {
// 	// Given a new Game
// 	g := state.NewGameState("game-id")

// 	// When we apply the WordBankStarted event
// 	g = g.Apply((state.WordAddedEvent{
// 		Word:     "apple",
// 		PlayerId: "player-1",
// 	}))

// 	// Then the phase is updated
// 	assert.Equal(t, g.SubmittedWords, []state.SubmittedWord{
// 		{
// 			Word:     "apple",
// 			PlayerId: "player-1",
// 		},
// 	})
// }

// func Test_Apply_Round1Started_UpdatesGamePhase(t *testing.T) {
// 	// Given a new Game
// 	g := state.NewGameState("game-id")

// 	// When we apply the WordBankStarted event
// 	g = g.Apply((state.Round1StartedEvent{}))

// 	// Then the phase is updated
// 	assert.Equal(t, g.Phase, "round-1")
// }

// func Test_Apply_Round2Started_UpdatesGamePhase(t *testing.T) {
// 	// Given a new Game
// 	g := state.NewGameState("game-id")

// 	// When we apply the WordBankStarted event
// 	g = g.Apply((state.Round2StartedEvent{}))

// 	// Then the phase is updated
// 	assert.Equal(t, g.Phase, "round-2")
// }

// func Test_Apply_Round3Started_UpdatesGamePhase(t *testing.T) {
// 	// Given a new Game
// 	g := state.NewGameState("game-id")

// 	// When we apply the WordBankStarted event
// 	g = g.Apply((state.Round3StartedEvent{}))

// 	// Then the phase is updated
// 	assert.Equal(t, g.Phase, "round-3")
// }
