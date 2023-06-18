package game

// Player represents one device connected to the game server, which
// sends and receives messages in real-time related to the game logic.
type Player struct {
	Id            string
	PlayerChannel PlayerChannel
}

func NewPlayer(id string, channel PlayerChannel) *Player {
	return &Player{
		Id:            id,
		PlayerChannel: channel,
	}
}
