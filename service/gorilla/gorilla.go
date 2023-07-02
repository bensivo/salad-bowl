package gorilla

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bensivo/salad-bowl/hub/adapters"
	"github.com/bensivo/salad-bowl/lobby"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartGorillaServer(svc *lobby.LobbyService) {
	r := mux.NewRouter()
	r.HandleFunc("/lobbies/{lobbyId}/connect", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		l, err := svc.GetLobby(vars["lobbyId"])
		if err != nil {
			writeJson(w, 500, map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Failed to get lobby: %v", err),
			})
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection", err)
			return
		}

		playerChannel := &adapters.WebsocketPlayerChannel{
			Conn: conn,
		}
		playerId := r.URL.Query().Get("playerId")

		if playerId != "" {
			l.Hub.HandleReconnection(playerChannel, playerId)
		} else {
			l.Hub.HandleNewConnection(playerChannel)
		}
	})

	r.HandleFunc("/lobbies", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request to create new lobby")

		lobbyId, err := svc.CreateNewLobby()
		if err != nil {
			fmt.Println("Error creating new lobby", err)
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}

		writeJson(w, 201, map[string]interface{}{
			"status":  "Success",
			"lobbyId": lobbyId,
		})
	}).
		Methods("POST")

	r.HandleFunc("/lobbies", func(w http.ResponseWriter, r *http.Request) {

		res := make(map[string]interface{})
		lobbies := svc.GetLobbies()

		for id, lobby := range lobbies {
			res[id] = lobby
		}

		writeJson(w, 201, res)
	}).
		Methods("GET")

	fmt.Println("Starting websocket server at port 8080")
	http.ListenAndServe(":8080", handlers.CORS()(r))
}

func writeJson(w http.ResponseWriter, statusCode int, payload map[string]interface{}) {
	resBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error serializing json response", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resBytes)
}
