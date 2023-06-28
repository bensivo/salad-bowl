package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bensivo/salad-bowl/gorilla"
	"github.com/bensivo/salad-bowl/hub"
	"github.com/bensivo/salad-bowl/lobby"
)

func main() {
	hub := hub.NewHub()
	lobby := lobby.NewLobby(hub)

	lobby.Start()

	go gorilla.StartGorillaServer(hub)

	waitForSigint()
}

func waitForSigint() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-c
		wg.Done()
	}()
	wg.Wait()
}
