package lobby

import (
	"errors"
	"fmt"
	"time"

	"github.com/bensivo/salad-bowl/hub"
)

type LobbyService struct {
	lobbies map[string]*Lobby
}

func NewLobbyService() *LobbyService {
	return &LobbyService{
		lobbies: make(map[string]*Lobby),
	}
}

func (svc *LobbyService) CreateNewLobby() (string, error) {
	if len(svc.lobbies) >= 10 {
		fmt.Printf("Could not create lobby. Instance full.\n")
		return "", errors.New("instance full. There are already 10 lobbies in the server")
	}

	newHub := hub.NewHub()

	newLobby := NewLobby(newHub)
	newLobby.Start()

	lobbyId := newLobby.ID
	svc.lobbies[lobbyId] = newLobby

	fmt.Printf("Created new lobby with id: %s\n", lobbyId)

	return lobbyId, nil
}

func (svc *LobbyService) DeleteLobby(id string) {
	fmt.Printf("Deleting lobby %s\n", id)
	delete(svc.lobbies, id)
}

func (svc *LobbyService) GetLobby(id string) (*Lobby, error) {
	l, ok := svc.lobbies[id]
	if !ok {
		return nil, fmt.Errorf("lobby %s not found", id)
	}

	return l, nil
}

func (svc *LobbyService) GetLobbies() map[string]*Lobby {
	return svc.lobbies
}

// Cleanup deletes all lobbies that have no players, and are more than 30 seconds old
func (svc *LobbyService) Cleanup() {
	now := time.Now()
lobbies:
	for id, lobby := range svc.lobbies {
		// Don't remove lobbies with at least 1 online player
		for _, player := range lobby.Players {
			if player.Status == "online" {
				continue lobbies
			}
		}

		// Don't remove brand new lobbies. Here, "brand new" was arbitrarily set to 30 seconds. May realize this needs to be longer in the future.
		if now.Sub(lobby.CreatedAt) < time.Duration(30*time.Second) {
			continue lobbies
		}

		fmt.Printf("Deleting empty lobby %s, created at %s\n", id, lobby.CreatedAt)
		delete(svc.lobbies, id)
	}
}

// StartCleanup schedules the cleanup job to run every 10 seconds
func (svc *LobbyService) StartCleanup() {
	go func() {
		for {
			time.Sleep(time.Duration(10 * time.Second))
			svc.Cleanup()
		}
	}()
}
