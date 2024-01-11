package state

import (
	"time"
)

type Player struct {
	PlayerID   string
	PlayerName string
	TeamName   string
}

type Team struct {
	TeamName string
	Score    int64
}

type SubmittedWord struct {
	Word     string `json:"word"`
	PlayerId string `json:"playerId"`
}

// GameState represents the state of a game at a single point in time
type GameState struct {
	ID               string          // Unique identifier for this game
	CreatedAt        time.Time       // When this game was created. Used for cleaning up old game instances
	Players          []Player        // List of players in the game
	Teams            []Team          // List of teams in the game
	Phase            string          // 'word-bank', 'round-1', 'round-2', 'round-3', or 'finished'
	SubmittedWords   []SubmittedWord // All words submitted during the word-bank phase
	RemainingWords   []string        // During a round, all the words left "in the bowl"
	RemainingPlayers []string        // During a round, all the players who have not played the charades role
}

func NewGameState(id string) GameState {
	return GameState{
		ID:               id,
		CreatedAt:        time.Now(),
		Players:          []Player{},
		Teams:            []Team{},
		Phase:            "",
		SubmittedWords:   []SubmittedWord{},
		RemainingWords:   []string{},
		RemainingPlayers: []string{},
	}
}

func (g GameState) Apply(e EventMetadata) GameState {
	// switch e.Name {
	// case PlayerJoined:
	// 	var event PlayerJoinedGameEvent = e.(PlayerJoinedEvent)

	// 	playerIndex := slices.IndexFunc(g.Players, func(p Player) bool {
	// 		return p.PlayerID == event.PlayerID
	// 	})
	// 	if playerIndex == -1 {
	// 		g.Players = append(g.Players, Player{
	// 			PlayerID:   event.PlayerID,
	// 			PlayerName: event.PlayerName,
	// 			TeamName:   "",
	// 		})
	// 	} else {
	// 		g.Players[playerIndex] = Player{
	// 			PlayerID:   event.PlayerID,
	// 			PlayerName: event.PlayerName,
	// 			TeamName:   g.Players[playerIndex].TeamName,
	// 		}
	// 	}

	// 	return g

	// case PlayerLeft:
	// 	var event PlayerLeftEvent = e.(PlayerLeftEvent)

	// 	g.Players = slices.DeleteFunc(g.Players, func(p Player) bool {
	// 		return p.PlayerID == event.PlayerID
	// 	})
	// 	return g

	// case TeamJoined:
	// 	var event TeamJoinedEvent = e.(TeamJoinedEvent)

	// 	for i, player := range g.Players {
	// 		if player.PlayerID == event.PlayerID {
	// 			g.Players[i].TeamName = event.TeamName
	// 		}
	// 	}
	// 	return g

	// case WordBankStarted:
	// 	g.Phase = "word-bank"
	// 	return g

	// case WordAdded:
	// 	var event WordAddedEvent = e.(WordAddedEvent)
	// 	g.SubmittedWords = append(g.SubmittedWords, SubmittedWord{
	// 		Word:     event.Word,
	// 		PlayerId: event.PlayerId,
	// 	})
	// 	return g

	// case Round1Started:
	// 	g.Phase = "round-1"
	// 	return g

	// case Round2Started:
	// 	g.Phase = "round-2"
	// 	return g

	// case Round3Started:
	// 	g.Phase = "round-3"
	// 	return g

	// case GuessingStarted:
	// 	g.RemainingWords = make([]string, len(g.SubmittedWords))
	// 	for i, word := range g.SubmittedWords {
	// 		g.RemainingWords[i] = word.Word
	// 	}

	// 	// TODO: randomly shuffle the words

	// 	return g

	// case GuessingFinished:
	// 	// TODO: notify all players, that the round is done?
	// 	// TODO: on ui show a loading screen, then wait for someone to press round2started?
	// 	return g

	// case GameFinished:
	// 	// TODO: show scores
	// 	return g
	// }

	// // You should never make it here, handle all possible events in the switch statement above
	panic("Unknown EventName")
}
