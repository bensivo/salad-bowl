package main

import (
	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/game/db"
	"github.com/bensivo/salad-bowl/service/pkg/game/http"
)

func main() {
	gameDb := db.NewInMemoryGameDb()
	gameSvc := game.NewGameService(gameDb)

	http.StartHttpGameService(gameSvc)
}
