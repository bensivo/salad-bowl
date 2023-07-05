package lobby

import (
	"fmt"
	"time"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/util"
)

type Lobby struct {
	ID        string
	Hub       hub.Hub
	Players   []*Player
	CreatedAt time.Time
}

func NewLobby(hub hub.Hub) *Lobby {
	lobbyId := util.RandStringId()

	teams := make([][]string, 2)
	teams[0] = []string{}
	teams[1] = []string{}

	return &Lobby{
		ID:        lobbyId,
		Hub:       hub,
		Players:   []*Player{},
		CreatedAt: time.Now(),
	}
}

// Start sets up all listeners and callbacks that the lobby uses during normal running.
func (l *Lobby) Start() {
	l.Hub.OnNewConnection(l.HandleNewConnection)
	l.Hub.OnMessage(l.HandleMessage)
	l.Hub.OnPlayerDisconnect(l.HandlePlayerDisconnect)
}

// HandleNewConnection adds a player to the lobby by id.
// It sends that player their welcome message, and then broadcasts the updated player list.
func (l *Lobby) HandleNewConnection(playerId string) {
	fmt.Printf("New connection with id %s\n", playerId)

	// Send the player's id to them so they know what it is
	welcomeMsg := hub.Message{
		Event: "notification.player-id",
		Payload: map[string]interface{}{
			"playerId": playerId,
		},
	}
	l.Hub.SendTo(playerId, welcomeMsg)

	isExistingPlayer := false
	for i := 0; i < len(l.Players); i++ {
		if l.Players[i].Id == playerId {
			l.Players[i].Status = "online"
			isExistingPlayer = true
			break
		}
	}

	if !isExistingPlayer {
		player := NewPlayer(playerId)
		l.Players = append(l.Players, player)
	}

	l.broadcastPlayerList()
}

func (l *Lobby) HandlePlayerDisconnect(playerId string) {
	fmt.Printf("Player %s disconnected\n", playerId)
	for i := 0; i < len(l.Players); i++ {
		player := l.Players[i]
		if player.Id == playerId {
			l.Players[i].Status = "offline"
		}
	}

	l.broadcastPlayerList()
}

func (l *Lobby) HandleMessage(playerId string, message hub.Message) {
	fmt.Printf("Received message from player %s: %v", playerId, message)
	switch message.Event {
	case "request.join-team":
		fmt.Printf("Player %s requesting to join team %d\n", playerId, message.Payload["team"])

		team := int(message.Payload["team"].(float64))
		requestId := message.Payload["requestId"].(string)

		if team != 0 && team != 1 {
			l.Hub.SendTo(playerId, hub.Message{
				Event: "response.join-team",
				Payload: map[string]interface{}{
					"requestId":   requestId,
					"status":      "error",
					"description": "Cannot join team. Please choose either team 0 or team 1",
					"team":        team,
				},
			})
			return
		}

		for i, player := range l.Players {
			if player.Id == playerId {
				l.Players[i].Team = team
				break
			}
		}
		// TODO: what if playerId not found, maybe we update the data struct to a map to help

		l.Hub.SendTo(playerId, hub.Message{
			Event: "response.join-team",
			Payload: map[string]interface{}{
				"requestId": requestId,
				"status":    "success",
				"team":      team,
			},
		})

		l.broadcastPlayerList()

	case "request.start-game":
		fmt.Printf("Player %s requesting to start the game\n", playerId)

		requestId := message.Payload["requestId"].(string)

		l.Hub.SendTo(playerId, hub.Message{
			Event: "response.start-game",
			Payload: map[string]interface{}{
				"requestId":   requestId,
				"status":      "success",
				"description": "Game starting",
			},
		})

		l.Hub.Broadcast(hub.Message{
			Event:   "notification.game-started",
			Payload: map[string]interface{}{},
		})
		return

	default:
		fmt.Printf("Unknown event %s\n", message.Event)
	}
}

func (l *Lobby) broadcastPlayerList() {
	players := []map[string]interface{}{}
	for _, player := range l.Players {
		players = append(players, map[string]interface{}{
			"id":     player.Id,
			"status": player.Status,
			"team":   player.Team,
		})
	}

	playerListMsg := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": players,
		},
	}
	l.Hub.Broadcast(playerListMsg)
}

func (l *Lobby) removeFromPlayers(playerId string) {
	for i := 0; i < len(l.Players); i++ {
		if l.Players[i].Id == playerId {
			l.Players[i] = l.Players[len(l.Players)-1]
			l.Players = l.Players[:len(l.Players)-1]
			break
		}
	}
}
