package game

import "time"

type Game struct {
	ID               string
	CreatedAt        time.Time       // When this game was created. Used for cleaning up old game instances
	Players          []Player        // List of players in the game
	Teams            []Team          // List of teams in the game
	Phase            string          // 'lobby', 'word-bank', 'round-1', 'round-2', 'round-3', or 'finished'
	SubmittedWords   []SubmittedWord // All words submitted during the word-bank phase
	RemainingWords   []string        // During a round, all the words left "in the bowl"
	RemainingPlayers []string        // During a round, all the players who have not played the charades role
}

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
