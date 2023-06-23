package game_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/game"
	"github.com/bensivo/salad-bowl/hub"
	"github.com/stretchr/testify/mock"
)

func TestGame_NewConnection_SendsPlayerId(t *testing.T) {
	mockHub := hub.NewMockHub(t)
	mockHub.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	mockHub.On("Broadcast", mock.Anything).Return(nil)

	// Given a new game
	game := game.NewGame(mockHub)

	// When a player is added
	game.AddPlayer("000-000")

	// Then the player receives an ID
	mockHub.AssertCalled(t, "SendTo", "000-000", map[string]interface{}{
		"ID": "000-000",
	})
}

func TestGame_NewConnection_BroadcastsPlayerList(t *testing.T) {
	mockHub := hub.NewMockHub(t)
	mockHub.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	mockHub.On("Broadcast", mock.Anything).Return(nil)

	// Given a new game
	game := game.NewGame(mockHub)

	// When 2 players are added
	game.AddPlayer("000-000")
	game.AddPlayer("111-111")

	// Then the player list is sent out
	mockHub.AssertCalled(t, "Broadcast", map[string]interface{}{
		"Players": []string{"000-000", "111-111"},
	})
}
