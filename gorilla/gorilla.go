package gorilla

import (
	"fmt"
	"net/http"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/hub/adapters"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartGorillaServer(h *hub.HubImpl) {
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

		h.HandleNewConnection(playerChannel)
	})

	fmt.Println("Starting websocket server at port 8080")
	http.ListenAndServe(":8080", nil)
}
