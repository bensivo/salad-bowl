package main

import (
	"github.com/bensivo/salad-bowl/service/pkg/game/db"
	"github.com/bensivo/salad-bowl/service/pkg/game/http"
)

func main() {
	gameService := db.NewInMemoryGameDb()

	http.StartHttpGameService(gameService)
}
