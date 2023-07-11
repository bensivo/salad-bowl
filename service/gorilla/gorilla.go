package gorilla

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bensivo/salad-bowl/game"
	"github.com/bensivo/salad-bowl/hub/adapters"
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

func StartGorillaServer(svc *game.GameService) {
	r := mux.NewRouter()
	r.HandleFunc("/game/{gameId}/connect", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		game, err := svc.GetOne(vars["gameId"])
		if err != nil {
			writeJson(w, 500, map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Failed to get game: %v", err),
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
			game.Hub.HandleReconnection(playerChannel, playerId)
		} else {
			game.Hub.HandleNewConnection(playerChannel)
		}
	})

	r.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request to create new game")

		gameId, err := svc.Create()
		if err != nil {
			fmt.Println("Error creating new game", err)
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}

		writeJson(w, 201, map[string]interface{}{
			"status": "Success",
			"gameId": gameId,
		})
	}).
		Methods("POST")

	r.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {

		res := make(map[string]interface{})
		games := svc.Get()

		for id, game := range games {
			res[id] = map[string]interface{}{
				"id":             game.ID,
				"phase":          game.Phase.Get(),
				"submittedWords": game.SubmittedWords.Get(),
				"players":        game.Players.Get(),
			}
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
