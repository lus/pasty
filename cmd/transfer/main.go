package main

import (
	"log"
	"os"

	"github.com/lus/pasty/internal/env"
	"github.com/lus/pasty/internal/storage"
)

func main() {
	// Validate the command line arguments
	if len(os.Args) != 3 {
		panic("Invalid command line arguments")
		return
	}

	// Load the optional .env file
	log.Println("Loading the optional .env file...")
	env.Load()

	// Create and initialize the first (from) driver
	from, err := storage.GetDriver(os.Args[1])
	if err != nil {
		panic(err)
	}
	err = from.Initialize()
	if err != nil {
		panic(err)
	}

	// Create and initialize the second (to) driver
	to, err := storage.GetDriver(os.Args[2])
	if err != nil {
		panic(err)
	}
	err = to.Initialize()
	if err != nil {
		panic(err)
	}

	// Retrieve a list of IDs from the first (from) driver
	ids, err := from.ListIDs()
	if err != nil {
		panic(err)
	}

	// Transfer every paste to the second (to) driver
	for _, id := range ids {
		log.Println("Transferring ID " + id + "...")

		// Retrieve the paste
		paste, err := from.Get(id)
		if err != nil {
			log.Println("[ERR]", err.Error())
			continue
		}

		// Save the paste
		err = to.Save(paste)
		if err != nil {
			log.Println("[ERR]", err.Error())
			continue
		}

		log.Println("Transferred ID " + id + ".")
	}
}
