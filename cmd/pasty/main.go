package main

import (
	"log"
	"time"

	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/storage"
	"github.com/lus/pasty/internal/web"
)

func main() {
	// Load the configuration
	log.Println("Loading the application configuration...")
	config.Load()

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
	if config.Current.AutoDelete.Enabled {
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
				time.Sleep(config.Current.AutoDelete.TaskInterval)
			}
		}()
	}

	// Serve the web resources
	log.Println("Serving the web resources...")
	panic(web.Serve())
}
