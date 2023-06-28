package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/gorilla/websocket"
)

type WebsocketPlayerChannel struct {
	Conn *websocket.Conn

	disconnectCb hub.DisconnectCallback
}

var _ hub.PlayerChannel = (*WebsocketPlayerChannel)(nil)

func (wpc *WebsocketPlayerChannel) Send(message interface{}) error {
	return wpc.Conn.WriteJSON(message)
}

func (wpc *WebsocketPlayerChannel) OnMessage(cb hub.MessageCallback) {
	go func() {
		defer func() {
			if wpc.disconnectCb != nil {
				wpc.disconnectCb()
			}
		}()

		for {
			_, bytes, err := wpc.Conn.ReadMessage()
			if err != nil {
				fmt.Println("Websocket closed:", err)
				return
			}

			var message hub.Message
			err = json.Unmarshal(bytes, &message)
			if err != nil {
				fmt.Printf("Failed to parse message to JSON, %v\n", err)
				continue
			}

			fmt.Printf("Received message: %v\n", message)
			cb(message)
		}
	}()
}

func (wpc *WebsocketPlayerChannel) OnDisconnect(cb hub.DisconnectCallback) {
	fmt.Println("Registering disconnect handler")
	wpc.disconnectCb = cb
}
