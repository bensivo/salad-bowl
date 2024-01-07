package db

import (
	"cmp"
	"errors"
	"slices"
	"time"

	"github.com/bensivo/salad-bowl/service/pkg/game"
)

var ErrNotFound = errors.New("not found")

type InMemoryGameDb struct {
	games map[string]*game.Game
}

var _ game.GameDb = (*InMemoryGameDb)(nil)

func NewInMemoryGameDb() game.GameDb {
	return &InMemoryGameDb{
		games: make(map[string]*game.Game),
	}
}

// Create implements GameDb.
func (gs *InMemoryGameDb) Create(ID string) (*game.Game, error) {
	game := &game.Game{
		ID:               ID,
		CreatedAt:        time.Now(),
		Players:          []game.Player{},
		Teams:            []game.Team{},
		Phase:            "lobby",
		SubmittedWords:   []game.SubmittedWord{},
		RemainingWords:   []string{},
		RemainingPlayers: []string{},
	}

	gs.games[ID] = game
	return game, nil
}

// GetAll implements GameDb.
func (gs *InMemoryGameDb) GetAll() ([]*game.Game, error) {
	games := make([]*game.Game, 0)

	for _, value := range gs.games {
		games = append(games, value)
	}

	// Golang doesn't guarantee ordering when iterating maps. To give consistent output, we sort by ID.
	slices.SortFunc(games, func(a *game.Game, b *game.Game) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return games, nil
}

// GetOne implements GameDb.
func (gs *InMemoryGameDb) GetOne(ID string) (*game.Game, error) {
	game, ok := gs.games[ID]
	if !ok {
		return nil, ErrNotFound
	}

	return game, nil
}

// Update implements GameDb.
func (gs *InMemoryGameDb) Update(ID string, game *game.Game) error {
	gs.games[ID] = game
	return nil
}

// Delete implements GameDb.
func (gs *InMemoryGameDb) Delete(ID string) error {
	delete(gs.games, ID)
	return nil
}
