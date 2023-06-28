package lobby

// Player represents one device connected to the game server, which
// sends and receives messages in real-time related to the game logic.
type Player struct {
	Id string
}

func NewPlayer(id string) *Player {
	return &Player{
		Id: id,
	}
}
