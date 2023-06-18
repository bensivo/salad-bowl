package game

import (
	"fmt"

	"github.com/bensivo/salad-bowl/util"
)

// Instance represents one game being played a group of players
type Instance struct {
	players []*Player

	game *Game
}

func NewInstance() *Instance {
	return &Instance{
		players: []*Player{},
		game:    NewGame(),
	}
}

func (i *Instance) HandleNewConnection(playerChannel PlayerChannel) {

	playerId := util.RandStringId()
	fmt.Printf("Creating new player: %s\n", playerId)

	player := NewPlayer(playerId, playerChannel)

	playerChannel.Send(map[string]interface{}{
		"ID": playerId,
	})

	i.AddPlayer(player)
}

func (i *Instance) AddPlayer(player *Player) {
	i.players = append(i.players, player)

	playerIds := make([]string, len(i.players))
	for i, player := range i.players {
		playerIds[i] = player.Id
	}

	i.Broadcast(map[string]interface{}{
		"players": playerIds,
	})
}

// Send a message to each player in the instance
func (i *Instance) Broadcast(message interface{}) error {

	for _, player := range i.players {

		fmt.Printf("Sending message to player %s: %v\n", player.Id, message)
		err := player.PlayerChannel.Send(message)

		if err != nil {
			fmt.Printf("Error sending message to player %s: %v\n", player.Id, err)
			return err
		}
	}

	return nil
}
