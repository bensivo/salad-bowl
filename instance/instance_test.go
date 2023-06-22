package instance_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/instance"
	"github.com/bensivo/salad-bowl/instance/adapters"
	"github.com/bensivo/salad-bowl/util"
	"github.com/stretchr/testify/assert"
)

func TestInstance_NewConnection_SendsPlayerId(t *testing.T) {
	util.SeedRand(0)

	// Given a new empty instance
	i := instance.NewInstance()

	// When 2 players connect
	pc1 := adapters.NewMockPlayerChannel()
	i.HandleNewConnection(pc1)
	pc2 := adapters.NewMockPlayerChannel()
	i.HandleNewConnection(pc2)

	// Then each player is added to the list
	assert.Equal(t, i.Players[0].Id, "SSN-9QH")
	assert.Equal(t, i.Players[1].Id, "RAM-LLM")

	// Then each player gets their id sent to them
	assert.Contains(t, pc1.Sent, map[string]interface{}{
		"ID": "SSN-9QH",
	})
	assert.Contains(t, pc2.Sent, map[string]interface{}{
		"ID": "RAM-LLM",
	})
}

func TestInstance_NewConnection_SendsPlayerList(t *testing.T) {
	util.SeedRand(0)

	// Given a new empty instance
	i := instance.NewInstance()

	// When 2 players connect
	pc1 := adapters.NewMockPlayerChannel()
	i.HandleNewConnection(pc1)
	pc2 := adapters.NewMockPlayerChannel()
	i.HandleNewConnection(pc2)

	// Each player receives the list of all players
	assert.Contains(t, pc1.Sent, map[string]interface{}{
		"players": []string{"SSN-9QH", "RAM-LLM"},
	})
	assert.Contains(t, pc2.Sent, map[string]interface{}{
		"players": []string{"SSN-9QH", "RAM-LLM"},
	})
}

func TestInstance_Broadcast_SendsToAllPlayers(t *testing.T) {
	// Given a new empty instance, with 2 players connected
	i := instance.NewInstance()
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
