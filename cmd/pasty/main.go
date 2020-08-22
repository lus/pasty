package main

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load the optional .env file
	log.Println("Loading the optional .env file...")
	env.Load()

	// Load the configured storage driver
	log.Println("Loading the configured storage driver...")
	err := storage.Load()
	if err != nil {
		panic(err)
	}

	// Wait for the program to exit
	// TODO: Replace this through blocking API server
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
