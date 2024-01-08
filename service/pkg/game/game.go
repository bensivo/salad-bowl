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
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

type Team struct {
	TeamName  string   `json:"teamName"`
	Score     int64    `json:"score"`
	PlayerIDs []string `json:"playerIds"`
}

type SubmittedWord struct {
	Word     string `json:"word"`
	PlayerId string `json:"playerId"`
}

// GameDb defines the functions used for persisting game state
type GameDb interface {
	Save(game *Game) error // acts like an 'upsert', inserts or updates based on game.ID
	GetAll() ([]*Game, error)
	GetOne(ID string) (*Game, error)
	Delete(ID string) error
}

// GameService defines the functions exposed by this module
type GameService interface {
	Create() (*Game, error)
	GetAll() ([]*Game, error)
	GetOne(ID string) (*Game, error)
	Delete(ID string) error

	HandleEvent(ID string, event GameEvent) error // Handle a received game event, applying it to the appropriate game, and making any updates to game state
	// TODO: how are people supposed to get notified of changes in game state? Do we add a subscriber pattern?
}
