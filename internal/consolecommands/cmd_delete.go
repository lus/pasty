package consolecommands

import (
	"context"
	"fmt"
)

func (router *Router) Delete(args []string) {
	if len(args) == 0 {
		fmt.Println("Expected 1 argument.")
		return
	}
	pasteID := args[0]
	paste, err := router.Storage.Pastes().FindByID(context.Background(), pasteID)
	if err != nil {
		fmt.Printf("Could not look up paste: %s.\n", err.Error())
		return
	}
	if paste == nil {
		fmt.Printf("Invalid paste ID: %s.\n", pasteID)
		return
	}
	if err := router.Storage.Pastes().DeleteByID(context.Background(), pasteID); err != nil {
		fmt.Printf("Could not delete paste: %s.\n", err.Error())
		return
	}
	fmt.Printf("Deleted paste %s.\n", pasteID)
}
