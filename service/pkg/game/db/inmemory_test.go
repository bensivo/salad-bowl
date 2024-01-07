package db_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bensivo/salad-bowl/service/pkg/game/db"
)

func assertNotNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		t.FailNow()
	}
}

func TestCreateReturnsGame(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	game, err := gameDb.Create("1111")
	assertNotNil(t, err)

	if game.ID != "1111" {
		t.Errorf("game.ID does not match passed ID")
		t.FailNow()
	}
}

func TestGetAllReturnsGames(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	// Add 3 games
	for i := 0; i < 3; i++ {
		_, err := gameDb.Create(fmt.Sprintf("%d", i))
		assertNotNil(t, err)
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

	_, err := gameDb.Create("1111")
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

func TestUpdateReplacesPreviousValue(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	g, err := gameDb.Create("1111")
	assertNotNil(t, err)

	g.Phase = "asdf"
	err = gameDb.Update("1111", g)
	assertNotNil(t, err)

	newG, err := gameDb.GetOne("1111")
	assertNotNil(t, err)

	if newG.Phase != "asdf" {
		t.Errorf("game object should be updated")
	}
}

func TestDeleteRemovesGame(t *testing.T) {
	gameDb := db.NewInMemoryGameDb()

	_, err := gameDb.Create("1111")
	assertNotNil(t, err)

	err = gameDb.Delete("1111")
	assertNotNil(t, err)

	g, _ := gameDb.GetOne("1111")
	if g != nil {
		t.Errorf("game was not deleted")
	}
}
