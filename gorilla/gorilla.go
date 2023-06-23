package gorilla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/hub/adapters"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartGorillaServer(hub *hub.HubImpl) {
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

		hub.HandleNewConnection(playerChannel)
	})

	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading HTTP request body", err)
			w.WriteHeader(400)
			return
		}

		payload := make(map[string]string)
		err = json.Unmarshal(bytes, &payload)
		if err != nil {
			fmt.Println("Error parsing HTTP request body", err)
			w.WriteHeader(400)
			return
		}

		hub.Broadcast(payload)
		w.WriteHeader(200)
	})

	http.HandleFunc("/sendTo", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading HTTP request body", err)
			w.WriteHeader(400)
			return
		}

		payload := make(map[string]interface{})
		err = json.Unmarshal(bytes, &payload)
		if err != nil {
			fmt.Println("Error parsing HTTP request body", err)
			w.WriteHeader(400)
			return
		}

		hub.SendTo(fmt.Sprintf("%v", payload["ID"]), payload["Message"])
		w.WriteHeader(200)
	})

	fmt.Println("Starting websocket server at port 8080")
	http.ListenAndServe(":8080", nil)
}
