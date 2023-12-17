package state

type Player struct {
	ID   string
	Name string
	Team string
}

type Team struct {
	Name      string
	Score     int64
	PlayerIds []string
}

type SubmittedWord struct {
	Word     string `json:"word"`
	PlayerId string `json:"playerId"`
}

// GameState represents the state of a game at a single point in time
type GameState struct {
	ID               string          // Unique identifier for this game
	Players          []Player        // List of players in the game
	Teams            []Team          // List of teams in the game
	Phase            string          // 'word-bank', 'round-1', 'round-2', 'round-3', or 'finished'
	SubmittedWords   []SubmittedWord // All words submitted during the word-bank phase
	RemainingWords   []string        // During a round, all the words left "in the bowl"
	RemainingPlayers []string        // During a round, all the players who have not played the charades role
}

var _ GameState = GameState{
	ID:    "asdf",
	Phase: "word-bank",
	Players: []Player{
		{
			ID:   "alice",
			Name: "alice",
		},
		{
			ID:   "bob",
			Name: "bob",
		},
		{
			ID:   "charlie",
			Name: "charlie",
		},
		{
			ID:   "david",
			Name: "david",
		},
	},
	Teams: []Team{
		{
			Name: "attackers",
			PlayerIds: []string{
				"alice",
				"bob",
			},
		},
		{
			Name: "defenders",
			PlayerIds: []string{
				"charlie",
				"david",
			},
		},
	},
	SubmittedWords: []SubmittedWord{
		{
			Word:     "apple",
			PlayerId: "alice",
		},
		{
			Word:     "banana",
			PlayerId: "bob",
		},
		{
			Word:     "carrot",
			PlayerId: "charlie",
		},
	},
	RemainingWords: []string{
		"apple",
		"banana",
		"carrot",
	},
}

/*
TODO: define all the events that can happen in a game

- player-joined
- player-left
- team-joined
- word-bank-started
- word-added
- round-1-started
- round-2-started
- round-3-started
- guessing-started
- word-guessed
- word-skipped
- guessing-finished
- game-finished

TODO: define fields for each piece, and how the game state changes in response
*/
