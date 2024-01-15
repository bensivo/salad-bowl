package listener

import (
	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/log"
	"github.com/gorilla/websocket"
)

// Implements game.GameListener using a gorilla websocket.
// Use with GameService.RegisterGameListener to send all game state changes over the given websocket connection
type WebSocketGameListener struct {
	Conn *websocket.Conn
}

var _ game.GameListener = (*WebSocketGameListener)(nil)

// OnChange implements game.GameListener.
func (wsgl *WebSocketGameListener) OnChange(g game.Game) {
	log.Infof("Sending game state to websocket %s \n", wsgl.Conn.LocalAddr())
	wsgl.Conn.WriteJSON(map[string]interface{}{
		"id":        g.ID,
		"createdAt": g.CreatedAt,
		"phase":     g.Phase,
		"players":   g.Players,
		"teams":     g.Teams,
	})
}
