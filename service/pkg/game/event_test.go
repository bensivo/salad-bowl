package game_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bensivo/salad-bowl/service/pkg/game"
)

func assertEqual(t *testing.T, actual any, expected any) {
	if actual != expected {
		t.Errorf("Value %v does not match expected %v\n", actual, expected)
		t.FailNow()
	}
}

func TestGameEventFromJSON(t *testing.T) {
	str := `{"name":"player-joined","timestamp":"2023-12-25T00:00:00Z00:00","playerId":"asdf","playerName":"John Doe"}`

	event := game.PlayerJoinedEvent{}
	err := json.Unmarshal([]byte(str), &event)
	if err != nil {
		fmt.Println(err)
	}

	assertEqual(t, event.Name, game.PlayerJoined)
	assertEqual(t, event.Timestamp, "2023-12-25T00:00:00Z00:00")
	assertEqual(t, event.PlayerID, "asdf")
	assertEqual(t, event.PlayerName, "John Doe")
}

func TestGameEventToJSON(t *testing.T) {
	event := game.PlayerJoinedEvent{
		EventMetadata: game.EventMetadata{
			Name:      game.PlayerJoined,
			Timestamp: "2023-12-25T00:00:00Z00:00",
		},
		PlayerID:   "asdf",
		PlayerName: "John Doe",
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
	}

	assertEqual(t, string(bytes), `{"name":"player-joined","timestamp":"2023-12-25T00:00:00Z00:00","playerId":"asdf","playerName":"John Doe"}`)
}
