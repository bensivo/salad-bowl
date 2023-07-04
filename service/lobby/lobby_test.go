package lobby_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/lobby"
	"github.com/stretchr/testify/mock"
)

func TestLobby_NewConnection_SendsPlayerId(t *testing.T) {
	mockHub := hub.NewMockHub(t)
	mockHub.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	mockHub.On("Broadcast", mock.Anything).Return(nil)

	// Given a new game
	l := lobby.NewLobby(mockHub)

	// When a player is added
	l.HandleNewConnection("player-id")

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

	// Given a new l
	l := lobby.NewLobby(mockHub)

	// When 2 players are added
	l.HandleNewConnection("000-000")
	l.HandleNewConnection("111-111")

	// Then the player list is sent out
	expected := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": []map[string]interface{}{
				{
					"id":     "000-000",
					"status": "online",
				},
				{
					"id":     "111-111",
					"status": "online",
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

	// Given a new l, with 2 players
	l := lobby.NewLobby(mockHub)
	l.HandleNewConnection("000-000")
	l.HandleNewConnection("111-111")

	// When a player disconnects
	l.HandlePlayerDisconnect("000-000")

	// Then the player list is sent out, and player 000-000 is now offline
	expected := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": []map[string]interface{}{
				{
					"id":     "000-000",
					"status": "offline",
				},
				{
					"id":     "111-111",
					"status": "online",
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
	l := lobby.NewLobby(h)
	l.Start()
	l.HandleNewConnection("000-000")

	// when that player sends a join team request
	l.HandleMessage("000-000", hub.Message{
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

	h := hub.NewMockHub(t)
	h.On("OnNewConnection", mock.Anything).Return()
	h.On("OnMessage", mock.Anything).Return()
	h.On("OnPlayerDisconnect", mock.Anything).Return()
	h.On("SendTo", mock.Anything, mock.Anything).Return(nil)
	h.On("Broadcast", mock.Anything).Return(nil)

	// given a new game, with 2 players
	l := lobby.NewLobby(h)
	l.Start()
	l.HandleNewConnection("000-000")
	l.HandleNewConnection("111-111")

	// each player joins a team
	l.HandleMessage("000-000", hub.Message{
		Event: "request.join-team",
		Payload: map[string]interface{}{
			"requestId": "00000000-0000-0000-0000-000000000000",
			"team":      float64(0), // JSON serialization makes everything a float64
		},
	})
	l.HandleMessage("111-111", hub.Message{
		Event: "request.join-team",
		Payload: map[string]interface{}{
			"requestId": "00000000-0000-0000-0000-000000000000",
			"team":      float64(1), // JSON serialization makes everything a float64
		},
	})

	// then the game sends the team list to everyone
	h.AssertCalled(t, "Broadcast", hub.Message{
		Event: "state.teams",
		Payload: map[string]interface{}{
			"teams": [][]string{
				{"000-000"},
				{"111-111"},
			},
		},
	})
}
