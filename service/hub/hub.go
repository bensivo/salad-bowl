package hub

import (
	"fmt"

	"github.com/bensivo/salad-bowl/util"
)

// Hub manages communication with player devices
//
// It exposes methods for sending messages to player devices, and allows users to register callbacks for received messages.
//
// In the context of this app, each individual game / lobby gets its own hub instance. This prevents players receiving messages from
// other games going on at the same time.
//
//go:generate mockery --name Hub
type Hub interface {
	SendTo(playerId string, message Message) error
	Broadcast(message Message) error

	OnNewConnection(cb NewConnectionCallback)
	OnPlayerDisconnect(cb PlayerDisconnectCallback)
	OnMessage(cb PlayerMessageCallback)

	HandleNewConnection(playerChannel PlayerChannel)                 // Add a player connection to the hub, generating a playerId
	HandleReconnection(playerChannel PlayerChannel, playerId string) // Add a player connection, using an existing playerId
}

type NewConnectionCallback func(playerId string)
type PlayerDisconnectCallback func(playerId string)
type PlayerMessageCallback func(playerId string, message Message)

type HubImpl struct {
	PlayerChannels        map[string]PlayerChannel
	newConnectionCallback NewConnectionCallback
	disconnectCallback    PlayerDisconnectCallback
	playerMessageCallback PlayerMessageCallback
}

var _ Hub = (*HubImpl)(nil)

func NewHub() *HubImpl {
	return &HubImpl{
		PlayerChannels: make(map[string]PlayerChannel),
	}
}

func (h *HubImpl) OnNewConnection(cb NewConnectionCallback) {
	h.newConnectionCallback = cb
}

func (h *HubImpl) OnPlayerDisconnect(cb PlayerDisconnectCallback) {
	h.disconnectCallback = cb
}

func (h *HubImpl) OnMessage(cb PlayerMessageCallback) {
	h.playerMessageCallback = cb
}

func (h *HubImpl) HandleReconnection(playerChannel PlayerChannel, playerId string) {
	fmt.Printf("Reconnection from player ID: %s\n", playerId)

	existingPlayerChannel, exists := h.PlayerChannels[playerId]
	if exists {
		fmt.Printf("Removing old channel for player ID: %s\n", playerId)
		err := existingPlayerChannel.Close()
		if err != nil {
			fmt.Printf("Error closing existing player channel %s: %v\n", playerId, err)
		}

		delete(h.PlayerChannels, playerId)
		h.disconnectCallback(playerId)
	}

	h.addPlayerConnection(playerChannel, playerId)
}

func (h *HubImpl) HandleNewConnection(playerChannel PlayerChannel) {
	playerId := util.RandStringId()
	fmt.Printf("New player connection. Assigning ID: %s\n", playerId)

	h.addPlayerConnection(playerChannel, playerId)
}

func (h *HubImpl) addPlayerConnection(playerChannel PlayerChannel, playerId string) {
	h.PlayerChannels[playerId] = playerChannel

	playerChannel.OnDisconnect(func() {
		fmt.Printf("Player connection %s disconnected\n", playerId)
		delete(h.PlayerChannels, playerId)
		h.disconnectCallback(playerId)
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
