package game

import "encoding/json"

type EventName string

const (
	PlayerJoined EventName = "player-joined" // Add a player to the game
	PlayerLeft   EventName = "player-left"   // Remove a player from the game
	TeamJoined   EventName = "team-joined"   // Join a team, or switch teams

	// WordBankStarted  EventName = "word-bank-started"
	// WordAdded        EventName = "word-added"
	// Round1Started    EventName = "round-1-Started"
	// Round2Started    EventName = "round-2-Started"
	// Round3Started    EventName = "round-3-Started"
	// GuessingStarted  EventName = "guessing-started"
	// WordGuessed      EventName = "word-guessed"
	// WordSkipped      EventName = "word-skipped"
	// GuessingFinished EventName = "guessing-finished"
	// GameFinished     EventName = "game-finished"
)

// Base object for all events. The type for 'Payload' will depend on which EventName is used.
type GameEvent struct {
	Name      EventName `json:"name"`
	Timestamp string    `json:"timestamp"`
	Payload   any       `json:"payload"`
}

type GameEventPayload interface {
	PlayerJoinedEventPayload | any
}

func ParseGameEventPayload[T GameEventPayload](payload any, ptr *T) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, ptr)
	if err != nil {
		return err
	}

	return nil
}

type PlayerJoinedEventPayload struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

type PlayerLeftEventPayload struct {
	PlayerID string `json:"playerId"`
}

type TeamJoinedEventPayload struct {
	PlayerID string
	TeamName string
}
