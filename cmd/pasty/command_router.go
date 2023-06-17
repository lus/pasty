package main

import (
	"bufio"
	"fmt"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/storage"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
	"syscall"
)

var whitespaceRegex = regexp.MustCompile("\\s+")

type consoleCommandRouter struct {
	Config  *config.Config
	Storage storage.Driver
}

func (router *consoleCommandRouter) Listen() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Err(err).Msg("Could not read console input.")
			continue
		}

		commandData := strings.Split(whitespaceRegex.ReplaceAllString(strings.TrimSpace(input), " "), " ")
		if len(commandData) == 0 {
			fmt.Println("Invalid command.")
			continue
		}

		handle := strings.ToLower(commandData[0])
		var args []string
		if len(commandData) > 1 {
			args = commandData[1:]
		}

		switch handle {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  help                     : Shows this overview")
			fmt.Println("  stop                     : Stops the application")
			fmt.Println("  setmodtoken <id> <token> : Changes the modification token of the paste with ID <id> to <token>")
			fmt.Println("  delete <id>              : Deletes the paste with ID <id>")
			fmt.Println("  cleanup <duration>       : Deletes all pastes that are older than <duration>")
			break
		case "stop":
			if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
				fmt.Printf("Could not send interrupt signal: %s.\nUse Ctrl+C instead.\n", err.Error())
				break
			}
			return
		case "setmodtoken":
			router.SetModificationToken(args)
			break
		case "delete":
			router.Delete(args)
			break
		case "cleanup":
			router.Cleanup(args)
			break
		default:
			fmt.Println("Invalid command.")
			break
		}
	}
}
