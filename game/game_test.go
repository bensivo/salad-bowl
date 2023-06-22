package game_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/game"
	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/hub/adapters"
	"github.com/bensivo/salad-bowl/util"
	"github.com/stretchr/testify/assert"
)

func TestGame_NewConnection_SendsPlayerIds(t *testing.T) {
	util.SeedRand(0)

	// Given a new game
	hub := hub.NewHub()
	game := game.NewGame(hub)

	game.Start()

	// When a new player channel connects to the hub
	pc := adapters.NewMockPlayerChannel()
	hub.HandleNewConnection(pc)

	// Then the player receives an ID
	assert.Contains(t, pc.Sent, map[string]interface{}{
		"ID": "SSN-9QH",
	})
}

func TestGame_NewConnection_BroadcastsPlayerList(t *testing.T) {
	util.SeedRand(0)

	// Given a new game
	hub := hub.NewHub()
	game := game.NewGame(hub)

	game.Start()

	// When a new player channel connects to the hub
	pc1 := adapters.NewMockPlayerChannel()
	hub.HandleNewConnection(pc1)
	pc2 := adapters.NewMockPlayerChannel()
	hub.HandleNewConnection(pc2)

	// Then the player receives an ID
	assert.Contains(t, pc1.Sent, map[string]interface{}{
		"Players": []string{"SSN-9QH", "RAM-LLM"},
	})
}
