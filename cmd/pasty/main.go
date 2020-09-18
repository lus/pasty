package main

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/storage"
	"github.com/Lukaesebrot/pasty/internal/web"
	"log"
	"time"
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
	defer func() {
		log.Println("Terminating the storage driver...")
		err := storage.Current.Terminate()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Schedule the AutoDelete task
	if env.Bool("AUTODELETE", false) {
		log.Println("Scheduling the AutoDelete task...")
		go func() {
			for {
				// Run the cleanup sequence
				deleted, err := storage.Current.Cleanup()
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("AutoDelete: Deleted %d expired pastes", deleted)

				// Wait until the process should repeat
				time.Sleep(env.Duration("AUTODELETE_TASK_INTERVAL", 5*time.Minute))
			}
		}()
	}

	// Serve the web resources
	log.Println("Serving the web resources...")
	panic(web.Serve())
}
