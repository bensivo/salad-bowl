package gorilla

import (
	"fmt"
	"net/http"

	"github.com/bensivo/salad-bowl/game"
	"github.com/bensivo/salad-bowl/game/adapters"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartGorillaServer(instance *game.Instance) {

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Incoming connection from %s\n", r.RemoteAddr)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection", err)
			return
		}

		playerChannel := &adapters.WebsocketPlayerChannel{
			Conn: conn,
		}

		instance.HandleNewConnection(playerChannel)
	})

	http.ListenAndServe(":8080", nil)
}
