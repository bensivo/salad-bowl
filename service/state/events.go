package state

type EventName string

const (
	PlayerJoined     EventName = "PlayerJoined"
	PlayerLeft       EventName = "PlayerLeft"
	TeamJoined       EventName = "TeamJoined"
	WordBankStarted  EventName = "WordBankStarted"
	WordAdded        EventName = "WordAdded"
	Round1Started    EventName = "Round1Started"
	Round2Started    EventName = "Round2Started"
	Round3Started    EventName = "Round3Started"
	GuessingStarted  EventName = "GuessingStarted"
	WordGuessed      EventName = "WordGuessed"
	WordSkipped      EventName = "WordSkipped"
	GuessingFinished EventName = "GuessingFinished"
	GameFinished     EventName = "GameFinished"
)

type Event interface {
	Name() EventName
}

// Payload for the PlayerJoined GameEvent
type PlayerJoinedEvent struct {
	PlayerID   string
	PlayerName string
}

func (e PlayerJoinedEvent) Name() EventName {
	return PlayerJoined
}

// Event for the PlayerLeft GameEvent
type PlayerLeftEvent struct {
	PlayerID string
}

func (e PlayerLeftEvent) Name() EventName {
	return PlayerLeft
}

// Event for the TeamJoined GameEvent
type TeamJoinedEvent struct {
	PlayerID string
	TeamName string
}

func (e TeamJoinedEvent) Name() EventName {
	return TeamJoined
}

// Event for the WordBankStarted GameEvent
type WordBankStartedEvent struct {
}

func (e WordBankStartedEvent) Name() EventName {
	return WordBankStarted
}

// Event for the WordAdded GameEvent
type WordAddedEvent struct {
	Word     string
	PlayerId string
}

func (e WordAddedEvent) Name() EventName {
	return WordAdded
}

// Event for the Round1Started GameEvent
type Round1StartedEvent struct {
}

func (e Round1StartedEvent) Name() EventName {
	return Round1Started
}

// Event for the Round2Started GameEvent
type Round2StartedEvent struct {
}

func (e Round2StartedEvent) Name() EventName {
	return Round2Started
}

// Event for the Round3Started GameEvent
type Round3StartedEvent struct {
}

func (e Round3StartedEvent) Name() EventName {
	return Round3Started
}

// Event for the GuessingStarted GameEvent
type GuessingStartedEvent struct {
}

func (e GuessingStartedEvent) Name() EventName {
	return GuessingStarted
}

// Event for the WordGuessed GameEvent
type WordGuessedEvent struct {
}

// Event for the WordSkipped GameEvent
type WordSkippedEvent struct {
}

// Event for the GuessingFinished GameEvent
type GuessingFinishedEvent struct {
}

// Event for the GameFinished GameEvent
type GameFinishedEvent struct {
}
