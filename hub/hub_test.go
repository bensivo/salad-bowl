package hub_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/hub/adapters"
	"github.com/bensivo/salad-bowl/util"
	"github.com/stretchr/testify/assert"
)

func TestInstance_NewConnection_Send(t *testing.T) {
	util.SeedRand(0) // Makes generated IDs deterministic

	// Given a hub
	h := hub.NewHub()

	// When a player connects
	pc1 := adapters.NewMockPlayerChannel()
	h.HandleNewConnection(pc1)

	// Then I can send messages to that player through the hub
	h.SendTo("SSN-9QH", map[string]interface{}{"Hello": "Word"})
	assert.Contains(t, pc1.Sent, map[string]interface{}{"Hello": "Word"})
}

func TestInstance_NewConnection_CallsCallback(t *testing.T) {
	// Given a registered connection callback
	called := false
	var callback hub.NewConnectionCallback = func(id string) {
		called = true
	}

	h := hub.NewHub()
	h.RegisterNewConnectionCallback(callback)

	// When a player connects
	pc1 := adapters.NewMockPlayerChannel()
	h.HandleNewConnection(pc1)

	// Then the callback is called
	assert.Equal(t, called, true)
}

func TestInstance_Broadcast_SendsToAllPlayers(t *testing.T) {
	// Given a new empty hub, with 2 players connected
	i := hub.NewHub()
	pc1 := adapters.NewMockPlayerChannel()
	pc2 := adapters.NewMockPlayerChannel()
	i.HandleNewConnection(pc1)
	i.HandleNewConnection(pc2)

	// When broadcast
	msg := map[string]interface{}{
		"hello": "world",
	}
	i.Broadcast(msg)

	// Each player receives the message
	assert.Contains(t, pc1.Sent, msg)
	assert.Contains(t, pc2.Sent, msg)
}
