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
	str := `{
		"name": "player-joined",
		"timestamp": "2023-12-25T00:00:00Z00:00",
		"payload": {
			"playerId": "asdf",
			"playerName": "John Doe"
		}
	}`

	var event game.GameEvent
	err := json.Unmarshal([]byte(str), &event)
	if err != nil {
		t.Errorf("failed parsing event %v\n", err)
	}

	assertEqual(t, event.Name, game.PlayerJoined)
	assertEqual(t, event.Timestamp, "2023-12-25T00:00:00Z00:00")

	var payload game.PlayerJoinedEventPayload
	err = game.ParseGameEventPayload(event.Payload, &payload)
	if err != nil {
		t.Errorf("failed parsing payload %v\n", err)
	}

	assertEqual(t, payload.PlayerID, "asdf")
	assertEqual(t, payload.PlayerName, "John Doe")
}

func TestGameEventToJSON(t *testing.T) {
	event := game.GameEvent{
		Name:      game.PlayerJoined,
		Timestamp: "2023-12-25T00:00:00Z00:00",
		Payload: game.PlayerJoinedEventPayload{
			PlayerID:   "asdf",
			PlayerName: "John Doe",
		},
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
	}

	assertEqual(t, string(bytes), `{"name":"player-joined","timestamp":"2023-12-25T00:00:00Z00:00","payload":{"playerId":"asdf","playerName":"John Doe"}}`)
}
