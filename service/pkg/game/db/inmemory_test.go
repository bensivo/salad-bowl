package db_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/game/db"
)

func assertNotNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		t.FailNow()
	}
}

func TestSave(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	err := gameDb.Save(&game.Game{
		ID: "1",
	})
	assertNotNil(t, err)
}

func TestGetAllReturnsGames(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	// Add 3 games
	for i := 0; i < 3; i++ {
		gameDb.Save(&game.Game{
			ID: fmt.Sprintf("%d", i),
		})
	}

	games, err := gameDb.GetAll()
	assertNotNil(t, err)

	for i := 0; i < 3; i++ {
		if games[i].ID != fmt.Sprintf("%d", i) {
			t.Errorf("game.ID does not match passed ID")
			t.FailNow()
		}
	}
}

func TestGetOneReturnsGame(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	err := gameDb.Save(&game.Game{
		ID: "1111",
	})
	assertNotNil(t, err)

	game, err := gameDb.GetOne("1111")
	assertNotNil(t, err)

	if game.ID != "1111" {
		t.Errorf("Returned a different gameID than expected")
		t.FailNow()
	}
}

func TestGetOneThrowsOnNotFound(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	g, err := gameDb.GetOne("1111")

	if g != nil {
		t.Errorf("game should be nil")
	}

	if !errors.Is(err, db.ErrNotFound) {
		t.Errorf("Expected error ErrNotFound. Got %v", err)
	}
}

func TestSaveReplacesPreviousValue(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	err := gameDb.Save(&game.Game{
		ID:    "1111",
		Phase: "fdsa",
	})
	assertNotNil(t, err)

	err = gameDb.Save(&game.Game{
		ID:    "1111",
		Phase: "asdf",
	})
	assertNotNil(t, err)

	newG, err := gameDb.GetOne("1111")
	assertNotNil(t, err)

	if newG.Phase != "asdf" {
		t.Errorf("game object should be updated")
	}
}

func TestDeleteRemovesGame(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	err := gameDb.Save(&game.Game{
		ID: "1111",
	})
	assertNotNil(t, err)

	err = gameDb.Delete("1111")
	assertNotNil(t, err)

	g, _ := gameDb.GetOne("1111")
	if g != nil {
		t.Errorf("game was not deleted")
	}
}
