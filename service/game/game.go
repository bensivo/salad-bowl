package game

import (
	"fmt"
	"time"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/util"
)

// SubmittedWord is a game-state DTO, represents a word submitted by a player to the word bank
type SubmittedWord struct {
	Word     string `json:"word"`
	PlayerId string `json:"playerId"`
}

// Game contains all the state related to a single instance of salad bowl being played
type Game struct {
	ID             string
	Hub            hub.Hub
	Players        []*Player
	CreatedAt      time.Time
	SubmittedWords []SubmittedWord
	Phase          string // lobby, word-bank, round1, round2, round3
}

func NewGame(hub hub.Hub) *Game {
	gameId := util.RandStringId()

	teams := make([][]string, 2)
	teams[0] = []string{}
	teams[1] = []string{}

	return &Game{
		ID:             gameId,
		Hub:            hub,
		Players:        []*Player{},
		CreatedAt:      time.Now(),
		SubmittedWords: []SubmittedWord{},
		Phase:          "lobby",
	}
}

// Start sets up all listeners and callbacks that the game uses during normal running.
func (g *Game) Start() {
	g.Hub.OnNewConnection(g.HandleNewConnection)
	g.Hub.OnMessage(g.HandleMessage)
	g.Hub.OnPlayerDisconnect(g.HandlePlayerDisconnect)
}

// HandleNewConnection adds a player to the game by id.
// It sends that player their welcome message, and then broadcasts the updated player list.
func (g *Game) HandleNewConnection(playerId string) {
	fmt.Printf("New connection with id %s\n", playerId)

	// Send the player's id to them so they know what it is
	welcomeMsg := hub.Message{
		Event: "notification.player-id",
		Payload: map[string]interface{}{
			"playerId": playerId,
		},
	}
	g.Hub.SendTo(playerId, welcomeMsg)

	isExistingPlayer := false
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Id == playerId {
			g.Players[i].Status = "online"
			isExistingPlayer = true
			break
		}
	}

	if !isExistingPlayer {
		player := NewPlayer(playerId)
		g.Players = append(g.Players, player)
	}

	g.broadcastPlayerList()

	g.Hub.SendTo(playerId, hub.Message{
		Event: "state.game-phase",
		Payload: map[string]interface{}{
			"phase": g.Phase,
		},
	})

	g.Hub.SendTo(playerId, hub.Message{
		Event: "state.word-bank",
		Payload: map[string]interface{}{
			"submittedWords": g.SubmittedWords,
		},
	})
}

func (g *Game) HandlePlayerDisconnect(playerId string) {
	fmt.Printf("Player %s disconnected\n", playerId)
	for i := 0; i < len(g.Players); i++ {
		player := g.Players[i]
		if player.Id == playerId {
			g.Players[i].Status = "offline"
		}
	}

	g.broadcastPlayerList()
}

func (g *Game) HandleMessage(playerId string, message hub.Message) {
	fmt.Printf("Received message from player %s: %v\n", playerId, message)
	switch message.Event {
	case "request.join-team":
		fmt.Printf("Player %s requesting to join team %d\n", playerId, message.Payload["team"])

		team := int(message.Payload["team"].(float64))
		requestId := message.Payload["requestId"].(string)

		if team != 0 && team != 1 {
			g.Hub.SendTo(playerId, hub.Message{
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

		for i, player := range g.Players {
			if player.Id == playerId {
				g.Players[i].Team = team
				break
			}
		}

		g.Hub.SendTo(playerId, hub.Message{
			Event: "response.join-team",
			Payload: map[string]interface{}{
				"requestId": requestId,
				"status":    "success",
				"team":      team,
			},
		})

		g.broadcastPlayerList()

	case "request.start-game":
		fmt.Printf("Player %s requesting to start the game\n", playerId)

		requestId := message.Payload["requestId"].(string)

		// TODO: handle these error cases:
		//  - There is only 1 player in the game
		g.Hub.SendTo(playerId, hub.Message{
			Event: "response.start-game",
			Payload: map[string]interface{}{
				"requestId": requestId,
				"status":    "success",
			},
		})

		g.Phase = "word-bank"
		g.Hub.Broadcast(hub.Message{
			Event: "state.game-phase",
			Payload: map[string]interface{}{
				"phase": g.Phase,
			},
		})
		return

	case "request.add-word":
		word := message.Payload["word"].(string)
		fmt.Printf("Player %s requesting to add word: %s\n", playerId, word)

		g.SubmittedWords = append(g.SubmittedWords, SubmittedWord{
			Word:     word,
			PlayerId: playerId,
		})

		// TODO: handle these error cases:
		//  - duplicate word
		//  - player already submitted 3 words
		requestId := message.Payload["requestId"].(string)
		g.Hub.SendTo(playerId, hub.Message{
			Event: "response.add-word",
			Payload: map[string]interface{}{
				"requestId": requestId,
				"status":    "success",
			},
		})

		g.Hub.Broadcast(hub.Message{
			Event: "state.word-bank",
			Payload: map[string]interface{}{
				"submittedWords": g.SubmittedWords,
			},
		})
		return

	default:
		fmt.Printf("Unknown event %s\n", message.Event)
	}
}

func (g *Game) broadcastPlayerList() {
	players := []map[string]interface{}{}
	for _, player := range g.Players {
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
	g.Hub.Broadcast(playerListMsg)
}
