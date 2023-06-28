package lobby

import (
	"fmt"
	"time"

	"github.com/bensivo/salad-bowl/hub"
)

type Lobby struct {
	Hub           hub.Hub
	Players       []*Player
	TeamPlayerIds [][]string
	CreatedAt     time.Time
}

func NewLobby(hub hub.Hub) *Lobby {
	teams := make([][]string, 2)
	teams[0] = []string{}
	teams[1] = []string{}

	return &Lobby{
		Hub:           hub,
		Players:       []*Player{},
		TeamPlayerIds: teams, // length 2 array of strings, representing playerIds
		CreatedAt:     time.Now(),
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
	fmt.Printf("New player with id %s\n", playerId)

	// Send the player's id to them so they know what it is
	welcomeMsg := hub.Message{
		Event: "notification.player-id",
		Payload: map[string]interface{}{
			"playerId": playerId,
		},
	}
	l.Hub.SendTo(playerId, welcomeMsg)

	player := NewPlayer(playerId)
	l.Players = append(l.Players, player)

	l.broadcastPlayerList()
	l.broadcastTeams()
}

func (l *Lobby) HandlePlayerDisconnect(playerId string) {
	fmt.Printf("Player %s disconnected\n", playerId)

	l.removeFromPlayers(playerId)
	l.removeFromTeams(playerId)

	l.broadcastPlayerList()
	l.broadcastTeams()
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

		l.removeFromTeams(playerId)

		l.TeamPlayerIds[team] = append(l.TeamPlayerIds[team], playerId)
		l.Hub.SendTo(playerId, hub.Message{
			Event: "response.join-team",
			Payload: map[string]interface{}{
				"requestId": requestId,
				"status":    "success",
				"team":      team,
			},
		})

		l.broadcastTeams()

	default:
		fmt.Printf("Unknown event %s\n", message.Event)
	}
}

func (l *Lobby) broadcastPlayerList() {
	playerList := make([]string, len(l.Players))
	for i := 0; i < len(l.Players); i++ {
		playerList[i] = l.Players[i].Id
	}

	playerListMsg := hub.Message{
		Event: "state.player-list",
		Payload: map[string]interface{}{
			"players": playerList,
		},
	}
	l.Hub.Broadcast(playerListMsg)
}

func (l *Lobby) broadcastTeams() {
	l.Hub.Broadcast(hub.Message{
		Event: "state.teams",
		Payload: map[string]interface{}{
			"teams": l.TeamPlayerIds,
		},
	})
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
func (l *Lobby) removeFromTeams(playerId string) {
	// Remove player from all teams, so we don't end up with duplicate entries
	team0 := l.TeamPlayerIds[0]
	team1 := l.TeamPlayerIds[1]
	for i, v := range team0 {
		if v == playerId {
			fmt.Printf("Removing player %s from team %d\n", playerId, 0)
			team0[i] = team0[len(team0)-1]
			team0 = team0[:len(team0)-1]
			l.TeamPlayerIds[0] = team0
		}
	}
	for i, v := range team1 {
		if v == playerId {
			fmt.Printf("Removing player %s from team %d\n", playerId, 1)
			team1[i] = team1[len(team1)-1]
			team1 = team1[:len(team1)-1]
			l.TeamPlayerIds[1] = team1
		}
	}
}
