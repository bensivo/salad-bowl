package hub_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInstance_NewConnection_CallsCallback(t *testing.T) {
	// Given a registered connection callback
	called := false
	var callback hub.NewConnectionCallback = func(id string) {
		called = true
	}

	h := hub.NewHub()
	h.RegisterNewConnectionCallback(callback)

	// When a player connects
	pc1 := hub.NewMockPlayerChannel(t)
	h.HandleNewConnection(pc1)

	// Then the callback is called
	assert.Equal(t, called, true)
}

func TestInstance_NewConnection_Send(t *testing.T) {
	util.SeedRand(0) // Makes generated IDs deterministic
	pc1 := hub.NewMockPlayerChannel(t)
	pc1.On("Send", mock.Anything).Return(nil)

	// Given a hub
	h := hub.NewHub()

	// When a player connects
	h.HandleNewConnection(pc1)

	// Then I can send messages to that player through the hub
	h.SendTo("SSN-9QH", map[string]interface{}{"Hello": "Word"})

	pc1.AssertCalled(t, "Send", map[string]interface{}{"Hello": "Word"})
}

func TestInstance_Broadcast_SendsToAllPlayers(t *testing.T) {
	pc1 := hub.NewMockPlayerChannel(t)
	pc1.On("Send", mock.Anything).Return(nil)
	pc2 := hub.NewMockPlayerChannel(t)
	pc2.On("Send", mock.Anything).Return(nil)

	// Given a new empty hub, with 2 players connected
	i := hub.NewHub()

	i.HandleNewConnection(pc1)
	i.HandleNewConnection(pc2)

	// When broadcast
	msg := map[string]interface{}{
		"hello": "world",
	}
	i.Broadcast(msg)

	// Each player receives the message
	pc1.AssertCalled(t, "Send", msg)
	pc2.AssertCalled(t, "Send", msg)
}
