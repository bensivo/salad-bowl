package adapters

import (
	"github.com/bensivo/salad-bowl/game"
	"github.com/gorilla/websocket"
)

type WebsocketPlayerChannel struct {
	Conn *websocket.Conn
}

var _ game.PlayerChannel = (*WebsocketPlayerChannel)(nil)

func (wpc *WebsocketPlayerChannel) Send(message interface{}) error {
	return wpc.Conn.WriteJSON(message)
}
func (wpc *WebsocketPlayerChannel) OnMessage(cb game.MessageCallback) {

}
