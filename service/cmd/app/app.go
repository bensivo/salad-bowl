package main

import (
	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/game/http"
)

func main() {
	gameService := game.NewGameService()

	http.StartHttpGameService(gameService)
}
