package main

import "github.com/Lukaesebrot/pasty/internal/env"

func main() {
	// Load the optional .env file
	env.Load()
}
