package game

import (
	"cmp"
	"errors"
	"slices"
	"time"
)

var ErrNotFound = errors.New("not found")

type GameService interface {
	Create(ID string) (*Game, error)
	GetAll() ([]*Game, error)
	GetOne(ID string) (*Game, error)
	Update(ID string, game *Game) error
	Delete(ID string) error
}

type gameService struct {
	games map[string]*Game
}

var _ GameService = (*gameService)(nil)

func NewGameService() GameService {
	return &gameService{
		games: make(map[string]*Game),
	}
}

// Create implements GameService.
func (gs *gameService) Create(ID string) (*Game, error) {
	game := &Game{
		ID:               ID,
		CreatedAt:        time.Now(),
		Players:          []Player{},
		Teams:            []Team{},
		Phase:            "lobby",
		SubmittedWords:   []SubmittedWord{},
		RemainingWords:   []string{},
		RemainingPlayers: []string{},
	}

	gs.games[ID] = game
	return game, nil
}

// GetAll implements GameService.
func (gs *gameService) GetAll() ([]*Game, error) {
	games := make([]*Game, 0)

	for _, value := range gs.games {
		games = append(games, value)
	}

	// Golang doesn't guarantee ordering when iterating maps. To give consistent output, we sort by ID.
	slices.SortFunc(games, func(a *Game, b *Game) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return games, nil
}

// GetOne implements GameService.
func (gs *gameService) GetOne(ID string) (*Game, error) {
	game, ok := gs.games[ID]
	if !ok {
		return nil, ErrNotFound
	}

	return game, nil
}

// Update implements GameService.
func (gs *gameService) Update(ID string, game *Game) error {
	gs.games[ID] = game
	return nil
}

// Delete implements GameService.
func (gs *gameService) Delete(ID string) error {
	delete(gs.games, ID)
	return nil
}
