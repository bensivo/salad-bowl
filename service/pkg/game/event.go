package game

import "encoding/json"

type EventName string

const (
	PlayerJoined EventName = "player-joined" // Add a new player to the game
	PlayerLeft   EventName = "player-left"   // Remove a player from the game
	// TeamJoined       EventName = "team-joined"
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

// EventMetadata contains meta information common to all events
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

// // Event for the TeamJoined GameEvent
// type TeamJoinedEvent struct {
// 	EventMetadata
// 	PlayerID string
// 	TeamName string
// }

// // Event for the WordBankStarted GameEvent
// type WordBankStartedEvent struct {
// 	EventMetadata
// }

// // Event for the WordAdded GameEvent
// type WordAddedEvent struct {
// 	EventMetadata
// 	Word     string
// 	PlayerId string
// }

// // Event for the Round1Started GameEvent
// type Round1StartedEvent struct {
// 	EventMetadata
// }

// // Event for the Round2Started GameEvent
// type Round2StartedEvent struct {
// 	EventMetadata
// }

// // Event for the Round3Started GameEvent
// type Round3StartedEvent struct {
// 	EventMetadata
// }

// // Event for the GuessingStarted GameEvent
// type GuessingStartedEvent struct {
// 	EventMetadata
// }

// // Event for the WordGuessed GameEvent
// type WordGuessedEvent struct {
// 	EventMetadata
// }

// // Event for the WordSkipped GameEvent
// type WordSkippedEvent struct {
// 	EventMetadata
// }

// // Event for the GuessingFinished GameEvent
// type GuessingFinishedEvent struct {
// 	EventMetadata
// }

// // Event for the GameFinished GameEvent
// type GameFinishedEvent struct {
// 	EventMetadata
// }
