package adapters

import (
	"github.com/bensivo/salad-bowl/hub"
	"github.com/gorilla/websocket"
)

type WebsocketPlayerChannel struct {
	Conn *websocket.Conn
}

var _ hub.PlayerChannel = (*WebsocketPlayerChannel)(nil)

func (wpc *WebsocketPlayerChannel) Send(message interface{}) error {
	return wpc.Conn.WriteJSON(message)
}
func (wpc *WebsocketPlayerChannel) OnMessage(cb hub.MessageCallback) {

}
