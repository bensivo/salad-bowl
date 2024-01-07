package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bensivo/salad-bowl/service/pkg/game"
	"github.com/bensivo/salad-bowl/service/pkg/game/db"
	"github.com/bensivo/salad-bowl/service/pkg/log"
	"github.com/bensivo/salad-bowl/service/pkg/util"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartHttpGameService(gameService game.GameDb) {
	r := mux.NewRouter()

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		logHTTP(r)
		log.Info("Creating new game")
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

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		logHTTP(r)
		log.Info("Getting all games")
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

	r.HandleFunc("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
		logHTTP(r)
		id, ok := mux.Vars(r)["id"]
		if !ok {
			writeErr(w, 400, errors.New("no id param given"))
			return
		}

		log.Infof("Getting game %s\n", id)
		g, err := gameService.GetOne(id)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				writeErr(w, 404, err)
				return
			}

			fmt.Println(err)
			writeErr(w, 500, err)
			return
		}

		writeJson(w, 200, map[string]interface{}{
			"id":             g.ID,
			"createdAt":      g.CreatedAt,
			"phase":          g.Phase,
			"submittedWords": g.SubmittedWords,
			"players":        g.Players,
		})
	}).
		Methods("GET")

	r.HandleFunc("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
		logHTTP(r)
		id, ok := mux.Vars(r)["id"]
		if !ok {
			writeErr(w, 400, errors.New("no id param given"))
			return
		}

		log.Infof("Deleting game %s\n", id)
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

	log.Info("Starting websocket server at port 8080")
	err := http.ListenAndServe(":8080", handlers.CORS()(r))
	if err != nil {
		log.Infof("Failed to start server: %v\n", err)
	}
}

func logHTTP(r *http.Request) {
	log.Infof("%s %s\n", r.Method, r.URL)
}

func writeJsonArr(w http.ResponseWriter, statusCode int, payload []map[string]interface{}) {
	resBytes, err := json.Marshal(payload)
	if err != nil {
		log.Infof("Error serializing json response %v\n", err)
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
		log.Infof("Error serializing json response %v\n", err)
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
