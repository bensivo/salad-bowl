package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bensivo/salad-bowl/game"
	"github.com/bensivo/salad-bowl/gorilla"
	"github.com/bensivo/salad-bowl/util"
)

func main() {
	util.Init()

	instance := game.NewInstance()

	go gorilla.StartGorillaServer(instance)

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
