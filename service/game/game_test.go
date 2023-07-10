package game_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/game"
	"github.com/bensivo/salad-bowl/hub"
	"github.com/stretchr/testify/mock"
)

func TestLobby_NewConnection_SendsPlayerId(t *testing.T) {
	mockHub := hub.NewMockHub(t)
	mockHub.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	mockHub.On("Broadcast", mock.Anything).Return(nil)

	// Given a new game
	g := game.NewGame(mockHub)

	// When a player is added
	g.HandleNewConnection("player-id")

	// Then the player receives an ID
	expected := hub.Message{
		Event: "notification.player-id",
		Payload: map[string]interface{}{
			"playerId": "player-id",
		},
	}
	mockHub.AssertCalled(t, "SendTo", "player-id", expected)
}

func TestLobby_NewConnection_BroadcastsPlayerList(t *testing.T) {
	mockHub := hub.NewMockHub(t)
	mockHub.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	mockHub.On("Broadcast", mock.Anything).Return(nil)

	// Given a new g
	g := game.NewGame(mockHub)

	// When 2 players are added
	g.HandleNewConnection("000-000")
	g.HandleNewConnection("111-111")

	// Then the player list is sent out
	expected := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": []map[string]interface{}{
				{
					"id":     "000-000",
					"status": "online",
					"team":   0,
				},
				{
					"id":     "111-111",
					"status": "online",
					"team":   0,
				},
			},
		},
	}
	mockHub.AssertCalled(t, "Broadcast", expected)
}

func TestLobby_Disconnect_BroadcastsPlayerOffline(t *testing.T) {
	mockHub := hub.NewMockHub(t)
	mockHub.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	mockHub.On("Broadcast", mock.Anything).Return(nil)

	// Given a new g, with 2 players
	g := game.NewGame(mockHub)
	g.HandleNewConnection("000-000")
	g.HandleNewConnection("111-111")

	// When a player disconnects
	g.HandlePlayerDisconnect("000-000")

	// Then the player list is sent out, and player 000-000 is now offline
	expected := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": []map[string]interface{}{
				{
					"id":     "000-000",
					"status": "offline",
					"team":   0,
				},
				{
					"id":     "111-111",
					"status": "online",
					"team":   0,
				},
			},
		},
	}
	mockHub.AssertCalled(t, "Broadcast", expected)
}

func TestLobby_TeamRequest_Success(t *testing.T) {

	h := hub.NewMockHub(t)
	h.On("OnNewConnection", mock.Anything).Return()
	h.On("OnMessage", mock.Anything).Return()
	h.On("OnPlayerDisconnect", mock.Anything).Return()
	h.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	h.On("Broadcast", mock.Anything).Return(nil)

	// given a new game, with one player
	g := game.NewGame(h)
	g.Start()
	g.HandleNewConnection("000-000")

	// when that player sends a join team request
	g.HandleMessage("000-000", hub.Message{
		Event: "request.join-team",
		Payload: map[string]interface{}{
			"requestId": "00000000-0000-0000-0000-000000000000",
			"team":      float64(1), // JSON serialization makes everything a float64
		},
	})

	// then the game sends a success response
	h.AssertCalled(t, "SendTo", "000-000", hub.Message{
		Event: "response.join-team",
		Payload: map[string]interface{}{
			"requestId": "00000000-0000-0000-0000-000000000000",
			"status":    "success",
			"team":      int(1), // We convert this to an int internally when we return it. It serializes to the same JSON.
		},
	})
}

func TestLobby_TeamRequest_StateUpdate(t *testing.T) {

	mockHub := hub.NewMockHub(t)
	mockHub.On("OnNewConnection", mock.Anything).Return()
	mockHub.On("OnMessage", mock.Anything).Return()
	mockHub.On("OnPlayerDisconnect", mock.Anything).Return()
	mockHub.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	mockHub.On("Broadcast", mock.Anything).Return(nil)

	// given a new game, with 2 players
	g := game.NewGame(mockHub)
	g.Start()
	g.HandleNewConnection("000-000")
	g.HandleNewConnection("111-111")

	// When each player joins a team
	g.HandleMessage("000-000", hub.Message{
		Event: "request.join-team",
		Payload: map[string]interface{}{
			"requestId": "00000000-0000-0000-0000-000000000000",
			"team":      float64(0), // JSON serialization makes everything a float64
		},
	})
	g.HandleMessage("111-111", hub.Message{
		Event: "request.join-team",
		Payload: map[string]interface{}{
			"requestId": "00000000-0000-0000-0000-000000000000",
			"team":      float64(1), // JSON serialization makes everything a float64
		},
	})

	// Then the player list is sent out, with players on appropriate teams
	expected := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": []map[string]interface{}{
				{
					"id":     "000-000",
					"status": "online",
					"team":   0,
				},
				{
					"id":     "111-111",
					"status": "online",
					"team":   1,
				},
			},
		},
	}

	mockHub.AssertCalled(t, "Broadcast", expected)
}
