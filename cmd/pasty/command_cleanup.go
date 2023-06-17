package main

import (
	"context"
	"fmt"
	"time"
)

func (router *consoleCommandRouter) Cleanup(args []string) {
	if len(args) == 0 {
		fmt.Println("Expected 1 argument.")
		return
	}
	lifetime, err := time.ParseDuration(args[0])
	if err != nil {
		fmt.Printf("Could not parse duration: %s.\n", err.Error())
		return
	}
	amount, err := router.Storage.Pastes().DeleteOlderThan(context.Background(), lifetime)
	if err != nil {
		if err != nil {
			fmt.Printf("Could not delete pastes: %s.\n", err.Error())
			return
		}
	}
	fmt.Printf("Deleted %d pastes older than %s.\n", amount, lifetime)
}
