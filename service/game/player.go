package game

// Player represents one device connected to the game server, which
// sends and receives messages in real-time related to the game logic.
type Player struct {
	Id     string `json:"id"`
	Status string `json:"status"` // online, offline, idle, away
	Team   int    `json:"team"`
}

func NewPlayer(id string) *Player {
	return &Player{
		Id:     id,
		Status: "online",
		Team:   0,
	}
}
