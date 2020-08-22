package main

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/storage"
	"github.com/Lukaesebrot/pasty/internal/web"
	"log"
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
	defer storage.Current.Terminate()

	// Serve the web server
	panic(web.Serve())
}
