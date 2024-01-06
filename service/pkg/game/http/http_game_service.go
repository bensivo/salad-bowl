package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/util"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartHttpGameService(gameService game.GameService) {
	r := mux.NewRouter()

	r.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request to create new game")

		game, err := gameService.Create(util.RandStringId())
		if err != nil {
			fmt.Println("Error creating new game", err)
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}

		writeJson(w, 201, map[string]interface{}{
			"id":             game.ID,
			"createdAt":      game.CreatedAt,
			"phase":          game.Phase,
			"submittedWords": game.SubmittedWords,
			"players":        game.Players,
		})
	}).
		Methods("POST")

	r.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		games, err := gameService.GetAll()
		if err != nil {
			fmt.Println(err)
			writeErr(w, 500, err)
			return
		}

		res := make([]map[string]interface{}, 0)
		for _, game := range games {
			res = append(res, map[string]interface{}{ // TODO: can we just annotate game with json field annotations?
				"id":             game.ID,
				"createdAt":      game.CreatedAt,
				"phase":          game.Phase,
				"submittedWords": game.SubmittedWords,
				"players":        game.Players,
			})
		}

		writeJsonArr(w, 201, res)
	}).
		Methods("GET")

	r.HandleFunc("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			writeErr(w, 400, errors.New("no id param given"))
			return
		}

		game, err := gameService.GetOne(id)
		if err != nil {
			fmt.Println(err)
			writeErr(w, 500, err)
			return
		}

		writeJson(w, 200, map[string]interface{}{
			"id":             game.ID,
			"createdAt":      game.CreatedAt,
			"phase":          game.Phase,
			"submittedWords": game.SubmittedWords,
			"players":        game.Players,
		})
	}).
		Methods("GET")

	r.HandleFunc("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			writeErr(w, 400, errors.New("no id param given"))
			return
		}

		err := gameService.Delete(id)
		if err != nil {
			fmt.Println(err)
			writeErr(w, 500, err)
			return
		}

		writeJson(w, 200, map[string]interface{}{
			"status": "success",
		})
	}).
		Methods("DELETE")

	fmt.Println("Starting websocket server at port 8080")
	http.ListenAndServe(":8080", handlers.CORS()(r))
}

func writeJsonArr(w http.ResponseWriter, statusCode int, payload []map[string]interface{}) {
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

func writeErr(w http.ResponseWriter, statusCode int, err error) {
	writeJson(w, statusCode, map[string]interface{}{
		"error": fmt.Sprintf("%v", err),
	})
}
