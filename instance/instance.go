package instance

import (
	"fmt"

	"github.com/bensivo/salad-bowl/game"
	"github.com/bensivo/salad-bowl/util"
)

// Instance represents one game being played a group of players
type Instance struct {
	Players []*Player

	Game *game.Game
}

func NewInstance() *Instance {
	return &Instance{
		Players: []*Player{},
		Game:    game.NewGame(),
	}
}

func (i *Instance) HandleNewConnection(playerChannel PlayerChannel) {
	playerId := util.RandStringId()
	fmt.Printf("Creating new player: %s\n", playerId)

	player := NewPlayer(playerId, playerChannel)
	i.Players = append(i.Players, player)

	playerChannel.Send(map[string]interface{}{
		"ID": playerId,
	})

	i.broadcastPlayerList(player)
}

// Send a message to each player in the instance
func (i *Instance) Broadcast(message interface{}) error {
	for _, player := range i.Players {

		fmt.Printf("Sending message to player %s: %v\n", player.Id, message)
		err := player.PlayerChannel.Send(message)

		if err != nil {
			fmt.Printf("Error sending message to player %s: %v\n", player.Id, err)
			return err
		}
	}

	return nil
}

func (i *Instance) broadcastPlayerList(player *Player) {
	playerIds := make([]string, len(i.Players))
	for i, player := range i.Players {
		playerIds[i] = player.Id
	}

	i.Broadcast(map[string]interface{}{
		"players": playerIds,
	})
}
