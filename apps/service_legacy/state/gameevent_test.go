package state_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bensivo/salad-bowl/state"
	"github.com/stretchr/testify/assert"
)

func Test_GameEvent_FromJSON(t *testing.T) {
	str := `{"name":"player-joined","timestamp":"2023-12-25T00:00:00Z00:00","playerId":"asdf","playerName":"John Doe"}`

	event := state.PlayerJoinedEvent{}
	err := json.Unmarshal([]byte(str), &event)
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, event.Name, state.PlayerJoined)
	assert.Equal(t, event.Timestamp, "2023-12-25T00:00:00Z00:00")
	assert.Equal(t, event.PlayerID, "asdf")
	assert.Equal(t, event.PlayerName, "John Doe")
}

func Test_GameEvent_ToJSON(t *testing.T) {
	event := state.PlayerJoinedEvent{
		EventMetadata: state.EventMetadata{
			Name:      state.PlayerJoined,
			Timestamp: "2023-12-25T00:00:00Z00:00",
		},
		PlayerID:   "asdf",
		PlayerName: "John Doe",
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, string(bytes), `{"name":"player-joined","timestamp":"2023-12-25T00:00:00Z00:00","playerId":"asdf","playerName":"John Doe"}`)
}
