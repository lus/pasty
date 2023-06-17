package consolecommands

import (
	"context"
	"fmt"
)

func (router *Router) SetModificationToken(args []string) {
	if len(args) < 2 {
		fmt.Println("Expected 2 arguments.")
		return
	}
	pasteID := args[0]
	newToken := args[1]
	paste, err := router.Storage.Pastes().FindByID(context.Background(), pasteID)
	if err != nil {
		fmt.Printf("Could not look up paste: %s.\n", err.Error())
		return
	}
	if paste == nil {
		fmt.Printf("Invalid paste ID: %s.\n", pasteID)
		return
	}
	paste.ModificationToken = newToken
	if err := paste.HashModificationToken(); err != nil {
		fmt.Printf("Could not hash modification token: %s.\n", err.Error())
		return
	}
	if err := router.Storage.Pastes().Upsert(context.Background(), paste); err != nil {
		fmt.Printf("Could not update paste: %s.\n", err.Error())
		return
	}
	fmt.Printf("Changed modification token of paste %s to %s.\n", pasteID, newToken)
}
