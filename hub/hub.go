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
	SendTo(playerId string, message interface{}) error
	Broadcast(message interface{}) error

	RegisterNewConnectionCallback(cb NewConnectionCallback)
}

type NewConnectionCallback func(playerId string)

var _ Hub = (*HubImpl)(nil)

// HubImpl manages all communication with player devices.
//
// It publishes connection IDs to registered listeners, and exposes methods for sending messages by ID
type HubImpl struct {
	PlayerChannels        map[string]PlayerChannel
	newConnectionCallback NewConnectionCallback
}

func NewHub() *HubImpl {
	return &HubImpl{
		PlayerChannels: make(map[string]PlayerChannel),
	}
}

func (h *HubImpl) RegisterNewConnectionCallback(cb NewConnectionCallback) {
	h.newConnectionCallback = cb
}

func (h *HubImpl) HandleNewConnection(playerChannel PlayerChannel) {
	playerId := util.RandStringId() // TODO: Provide a mechanism for the channel to reuse a previously sent ID - for reconnects
	fmt.Printf("New player connection. Assigning ID %s\n", playerId)

	h.PlayerChannels[playerId] = playerChannel

	if h.newConnectionCallback != nil {
		h.newConnectionCallback(playerId)
	}
}

// SendTo sends a message to a single player channel by Id
func (h *HubImpl) SendTo(playerId string, message interface{}) error {
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
func (h *HubImpl) Broadcast(message interface{}) error {
	fmt.Println("Broadcasting message: ", message)
	for id := range h.PlayerChannels {
		h.SendTo(id, message)
	}

	return nil
}
