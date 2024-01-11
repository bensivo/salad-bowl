package game

import (
	"fmt"
	"time"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/observable"
	"github.com/bensivo/salad-bowl/util"
)

// SubmittedWord is a game-state DTO, represents a word submitted by a player to the word bank
type SubmittedWord struct {
	Word     string `json:"word"`
	PlayerId string `json:"playerId"`
}

// Game contains all the state related to a single instance of salad bowl being played
type Game struct {
	ID        string
	Hub       hub.Hub
	CreatedAt time.Time

	// Observable state objects. Any updates to these objects get broadcast to all players in the game
	Players        observable.Observable[[]Player]
	SubmittedWords observable.Observable[[]SubmittedWord]
	Phase          observable.Observable[string]
}

func NewGame(hub hub.Hub) *Game {
	gameId := util.RandStringId()

	teams := make([][]string, 2)
	teams[0] = []string{}
	teams[1] = []string{}

	players := observable.New([]Player{})
	submittedWords := observable.New([]SubmittedWord{})
	phase := observable.New("lobby")

	return &Game{
		ID:        gameId,
		Hub:       hub,
		CreatedAt: time.Now(),

		Players:        players,
		SubmittedWords: submittedWords,
		Phase:          phase,
	}
}

// Start sets up all listeners and callbacks that the game uses during normal running.
func (g *Game) Start() {
	g.Hub.OnNewConnection(g.HandleNewConnection)
	g.Hub.OnMessage(g.HandleMessage)
	g.Hub.OnPlayerDisconnect(g.HandlePlayerDisconnect)

	g.Players.OnChange(func(players []Player) {
		playersAsMap := []map[string]interface{}{}
		for _, player := range players {
			playersAsMap = append(playersAsMap, map[string]interface{}{
				"id":     player.Id,
				"status": player.Status,
				"team":   player.Team,
			})
		}
		g.Hub.Broadcast(hub.Message{
			Event: "state.player-list",
			Payload: map[string]interface{}{
				"players": playersAsMap,
			},
		})
	})

	g.Phase.OnChange(func(value string) {
		g.Hub.Broadcast(hub.Message{
			Event: "state.game-phase",
			Payload: map[string]interface{}{
				"phase": value,
			},
		})
	})

	g.SubmittedWords.OnChange(func(value []SubmittedWord) {
		g.Hub.Broadcast(hub.Message{
			Event: "state.word-bank",
			Payload: map[string]interface{}{
				"submittedWords": value,
			},
		})
	})
}

// HandleNewConnection adds a player to the game by id.
// It sends that player their welcome message, and then sends them all the observable state values they need
func (g *Game) HandleNewConnection(playerId string) {
	fmt.Printf("New connection with id %s\n", playerId)

	g.Hub.SendTo(playerId, hub.Message{
		Event: "notification.player-id",
		Payload: map[string]interface{}{
			"playerId": playerId,
		},
	})

	players := g.Players.Get()
	isExistingPlayer := false

	for i := 0; i < len(players); i++ {
		if players[i].Id == playerId {
			players[i].Status = "online"
			isExistingPlayer = true
			break
		}
	}

	if !isExistingPlayer {
		player := Player{
			Id:     playerId,
			Status: "online",
			Team:   0,
		}
		players = append(players, player)
	}

	g.Players.Set(players)

	g.Hub.SendTo(playerId, hub.Message{
		Event: "state.game-phase",
		Payload: map[string]interface{}{
			"phase": g.Phase.Get(),
		},
	})

	g.Hub.SendTo(playerId, hub.Message{
		Event: "state.word-bank",
		Payload: map[string]interface{}{
			"submittedWords": g.SubmittedWords.Get(),
		},
	})
}

func (g *Game) HandlePlayerDisconnect(playerId string) {
	fmt.Printf("Player %s disconnected\n", playerId)
	players := g.Players.Get()
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.Id == playerId {
			players[i].Status = "offline"
		}
	}

	g.Players.Set(players)
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

		players := g.Players.Get()
		for i, player := range players {
			if player.Id == playerId {
				players[i].Team = team
				break
			}
		}
		g.Players.Set(players)

		g.Hub.SendTo(playerId, hub.Message{
			Event: "response.join-team",
			Payload: map[string]interface{}{
				"requestId": requestId,
				"status":    "success",
				"team":      team,
			},
		})

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

		g.Phase.Set("word-bank")
		return

	case "request.add-word":
		requestId := message.Payload["requestId"].(string)

		word := message.Payload["word"].(string)
		fmt.Printf("Player %s requesting to add word: %s\n", playerId, word)

		words := g.SubmittedWords.Get()

		for _, submittedWord := range words {
			if word == submittedWord.Word {
				fmt.Printf("Cannot add word %s. Duplicate.\n", word)

				g.Hub.SendTo(playerId, hub.Message{
					Event: "response.add-word",
					Payload: map[string]interface{}{
						"requestId":   requestId,
						"status":      "error",
						"description": fmt.Sprintf("\"%s\" is already in the word bank", word),
					},
				})
				return
			}
		}

		words = append(words, SubmittedWord{
			Word:     word,
			PlayerId: playerId,
		})
		g.SubmittedWords.Set(words)

		// TODO: handle these error cases:
		//  - duplicate word
		//  - player already submitted 3 words
		g.Hub.SendTo(playerId, hub.Message{
			Event: "response.add-word",
			Payload: map[string]interface{}{
				"requestId": requestId,
				"status":    "success",
			},
		})

		return

	default:
		fmt.Printf("Unknown event %s\n", message.Event)
	}
}
