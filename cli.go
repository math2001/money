package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type App struct {
	cryptor  *Cryptor
	commands map[string]func(...string) error
}

func (app *App) Start() {
	app.commands = map[string]func(...string) error{
		"help": app.help,
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		instructionline, err := reader.ReadString('\n')
		if err != nil {
			// FIXME: retry a few times, and quit if it doesn't work
			log.Fatalf("reading from stdin: %s", err)
		}
		command, args := parse(instructionline)

		if command == "exit" {
			return
		}
		fn, ok := app.commands[command]
		if !ok {
			app.unhandledCommand(command, args)
		} else {
			err := fn(args...)
			if err != nil {
				fmt.Printf("[err] %s\n", err)
			}
		}
	}
}

func (app *App) help(args ...string) error {
	fmt.Println("This is the help message!")
	if len(args) > 0 {
		return fmt.Errorf("doesn't take any argument yet")
	}
	return nil
}

func (app *App) unhandledCommand(command string, args []string) {
	// FIXME: look up similar commands
	fmt.Printf("Command %q doesn't exist\n", command)
}

// FIXME: this might need to be fancier (support quotes and stuff).
func parse(instructionline string) (string, []string) {
	fields := strings.Fields(instructionline)
	return fields[0], fields[1:]
}

func NewApp(cryptor *Cryptor) *App {
	return &App{
		cryptor: cryptor,
	}
}
