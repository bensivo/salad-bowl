package state

type EventName string

const (
	PlayerJoined     EventName = "player-joined"
	PlayerLeft       EventName = "player-left"
	TeamJoined       EventName = "team-joined"
	WordBankStarted  EventName = "word-bank-started"
	WordAdded        EventName = "word-added"
	Round1Started    EventName = "round-1-Started"
	Round2Started    EventName = "round-2-Started"
	Round3Started    EventName = "round-3-Started"
	GuessingStarted  EventName = "guessing-started"
	WordGuessed      EventName = "word-guessed"
	WordSkipped      EventName = "word-skipped"
	GuessingFinished EventName = "guessing-finished"
	GameFinished     EventName = "game-finished"
)

// EventMetadata contains meta information common to all events
type EventMetadata struct {
	Name      EventName `json:"name"`
	Timestamp string    `json:"timestamp"`
}

// Payload for the PlayerJoined GameEvent
type PlayerJoinedEvent struct {
	EventMetadata
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

// Event for the PlayerLeft GameEvent
type PlayerLeftEvent struct {
	EventMetadata
	PlayerID string
}

// Event for the TeamJoined GameEvent
type TeamJoinedEvent struct {
	EventMetadata
	PlayerID string
	TeamName string
}

// Event for the WordBankStarted GameEvent
type WordBankStartedEvent struct {
	EventMetadata
}

// Event for the WordAdded GameEvent
type WordAddedEvent struct {
	EventMetadata
	Word     string
	PlayerId string
}

// Event for the Round1Started GameEvent
type Round1StartedEvent struct {
	EventMetadata
}

// Event for the Round2Started GameEvent
type Round2StartedEvent struct {
	EventMetadata
}

// Event for the Round3Started GameEvent
type Round3StartedEvent struct {
	EventMetadata
}

// Event for the GuessingStarted GameEvent
type GuessingStartedEvent struct {
	EventMetadata
}

// Event for the WordGuessed GameEvent
type WordGuessedEvent struct {
	EventMetadata
}

// Event for the WordSkipped GameEvent
type WordSkippedEvent struct {
	EventMetadata
}

// Event for the GuessingFinished GameEvent
type GuessingFinishedEvent struct {
	EventMetadata
}

// Event for the GameFinished GameEvent
type GameFinishedEvent struct {
	EventMetadata
}
