package game

import "time"

// Game defines all the state needed for a given game instance
type Game struct {
	ID        string    `json:"id"`        // Unique Game Id, usually follows the form "XXX-XXX"
	CreatedAt time.Time `json:"createdAt"` // ISO8601 timestamp of when this game was created
	Players   []Player  `json:"players"`   // List of players in the game
	Teams     []Team    `json:"teams"`     // List of teams in the game, by default each game has 3 teams. "Red", "Blue", and "Spectators"
	Phase     string    `json:"phase"`     // 'lobby', 'word-bank', 'round-1', 'round-2', 'round-3', or 'finished'
	// SubmittedWords   []SubmittedWord `json:"submittedWords"`   // All words submitted during the word-bank phase
	// RemainingWords   []string        `json:"remainingWords"`   // During a round, all the words left "in the bowl"
	// RemainingPlayers []string        `json:"remainingPlayers"` // During a round, all the players who have not played the charades role
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

// type SubmittedWord struct {
// 	Word     string `json:"word"`
// 	PlayerId string `json:"playerId"`
// }

// GameDb defines the functions used for persisting game state
type GameDb interface {
	Save(game *Game) error // acts like an 'upsert', inserts or updates based on game.ID
	GetAll() ([]*Game, error)
	GetOne(ID string) (*Game, error)
	Delete(ID string) error
}

// GameService defines the functions exposed by this module
type GameService interface {
	// CRD functions for game instances
	Create() (*Game, error)
	GetAll() ([]*Game, error)
	GetOne(ID string) (*Game, error)
	Delete(ID string) error

	// Handle a received game event, applying it to the appropriate game, and making any updates to game state
	HandleEvent(ID string, event GameEvent) error

	// Register a listener, which shoudl be called anytime the game state changes
	RegisterListener(ID string, listener GameListener) error

	// Remove a previously registered listener
	DeregisterListener(ID string, listener GameListener) error
}

// Used for notifying clients of changes in game state in real time.
// Could be implemented using websockets, or SSEs, or even something like MQTT
type GameListener interface {
	OnChange(g Game)

	// TODO: include a close() method, for when the game is deleted and the listener should stop
}
