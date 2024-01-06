package game_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bensivo/salad-bowl/service/pkg/game"
)

func assertNotNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		t.FailNow()
	}
}

func TestCreateReturnsGame(t *testing.T) {
	gameService := game.NewGameService()

	game, err := gameService.Create("1111")
	assertNotNil(t, err)

	if game.ID != "1111" {
		t.Errorf("game.ID does not match passed ID")
		t.FailNow()
	}
}

func TestGetAllReturnsGames(t *testing.T) {
	gameService := game.NewGameService()

	// Add 3 games
	for i := 0; i < 3; i++ {
		_, err := gameService.Create(fmt.Sprintf("%d", i))
		assertNotNil(t, err)
	}

	games, err := gameService.GetAll()
	assertNotNil(t, err)

	for i := 0; i < 3; i++ {
		if games[i].ID != fmt.Sprintf("%d", i) {
			t.Errorf("game.ID does not match passed ID")
			t.FailNow()
		}
	}
}

func TestGetOneReturnsGame(t *testing.T) {
	gameService := game.NewGameService()

	_, err := gameService.Create("1111")
	assertNotNil(t, err)

	game, err := gameService.GetOne("1111")
	assertNotNil(t, err)

	if game.ID != "1111" {
		t.Errorf("Returned a different gameID than expected")
		t.FailNow()
	}
}

func TestGetOneThrowsOnNotFound(t *testing.T) {
	gameService := game.NewGameService()

	g, err := gameService.GetOne("1111")

	if g != nil {
		t.Errorf("game should be nil")
	}

	if !errors.Is(err, game.ErrNotFound) {
		t.Errorf("Expected error ErrNotFound. Got %v", err)
	}
}

func TestUpdateReplacesPreviousValue(t *testing.T) {
	gameService := game.NewGameService()

	g, err := gameService.Create("1111")
	assertNotNil(t, err)

	g.Phase = "asdf"
	err = gameService.Update("1111", g)
	assertNotNil(t, err)

	newG, err := gameService.GetOne("1111")
	assertNotNil(t, err)

	if newG.Phase != "asdf" {
		t.Errorf("game object should be updated")
	}
}

func TestDeleteRemovesGame(t *testing.T) {
	gameService := game.NewGameService()

	_, err := gameService.Create("1111")
	assertNotNil(t, err)

	err = gameService.Delete("1111")
	assertNotNil(t, err)

	g, _ := gameService.GetOne("1111")
	if g != nil {
		t.Errorf("game was not deleted")
	}
}
