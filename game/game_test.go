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
	expected := hub.Message{
		Event: "notification.player-id",
		Payload: map[string]interface{}{
			"playerId": "000-000",
		},
	}
	mockHub.AssertCalled(t, "SendTo", "000-000", expected)
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
	expected := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": []string{"000-000", "111-111"},
		},
	}
	mockHub.AssertCalled(t, "Broadcast", expected)
}
