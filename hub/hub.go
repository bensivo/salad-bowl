package hub

import (
	"fmt"

	"github.com/bensivo/salad-bowl/util"
)

// Hub manages communication with player devices
//
// It exposes methods for sending messages to player devices, and allows users to register callbacks for received messages.
//
//go:generate mockery --name Hub
type Hub interface {
	SendTo(playerId string, message Message) error
	Broadcast(message Message) error

	OnNewConnection(cb NewConnectionCallback)
	OnMessage(cb PlayerMessageCallback)
}

type NewConnectionCallback func(playerId string)
type PlayerMessageCallback func(playerId string, message Message)

var _ Hub = (*HubImpl)(nil)

// HubImpl manages all communication with player devices.
//
// It publishes connection IDs to registered listeners, and exposes methods for sending messages by ID
type HubImpl struct {
	PlayerChannels        map[string]PlayerChannel
	newConnectionCallback NewConnectionCallback
	playerMessageCallback PlayerMessageCallback
}

func NewHub() *HubImpl {
	return &HubImpl{
		PlayerChannels: make(map[string]PlayerChannel),
	}
}

func (h *HubImpl) OnNewConnection(cb NewConnectionCallback) {
	h.newConnectionCallback = cb
}

func (h *HubImpl) OnMessage(cb PlayerMessageCallback) {
	h.playerMessageCallback = cb
}

func (h *HubImpl) HandleNewConnection(playerChannel PlayerChannel) {
	playerId := util.RandStringId()
	fmt.Printf("New player connection. Assigning ID %s\n", playerId)

	h.PlayerChannels[playerId] = playerChannel

	playerChannel.OnDisconnect(func() {
		fmt.Printf("Player connection %s disconnected\n", playerId)
		delete(h.PlayerChannels, playerId)
	})

	if h.newConnectionCallback != nil {
		h.newConnectionCallback(playerId)
	}

	playerChannel.OnMessage(func(message Message) {
		if h.playerMessageCallback != nil {
			h.playerMessageCallback(playerId, message)
		}
	})
}

// SendTo sends a message to a single player channel by Id
func (h *HubImpl) SendTo(playerId string, message Message) error {
	fmt.Printf("Sending message to %s: %v\n", playerId, message)
	pc, exists := h.PlayerChannels[playerId]
	if !exists {
		return fmt.Errorf("player %s not found", playerId)
	}

	err := pc.Send(message)
	if err != nil {
		fmt.Printf("Failed sending message to player %s\n", playerId)
		return err
	}
	return nil
}

// Broadcasts sends a message to all players
func (h *HubImpl) Broadcast(message Message) error {
	fmt.Println("Broadcasting message: ", message)
	for id := range h.PlayerChannels {
		h.SendTo(id, message)
	}

	return nil
}
