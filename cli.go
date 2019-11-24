package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Cli struct {
	cryptor  *Cryptor
	commands map[string]func(...string) error
}

func (cli *Cli) Start() {
	cli.commands = map[string]func(...string) error{
		"help": cli.help,
		"ls":   cli.ls,
		"load": cli.load,
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
		fn, ok := cli.commands[command]
		if !ok {
			cli.unhandledCommand(command, args)
		} else {
			err := fn(args...)
			if err != nil {
				fmt.Printf("[err] %s\n", err)
			}
		}
	}
}

func (cli *Cli) load(args ...string) error {
	if len(args) != 1 {
		// FIXME: display usage for this command
		return fmt.Errorf("load takes one argument, the filename")
	}
	path := filepath.Join(store, args[0])
	content, err := cli.cryptor.Load(path)
	if err != nil {
		return fmt.Errorf("loading %q: %s", path, err)
	}
	fmt.Println(string(content))
	return nil
}

func (cli *Cli) ls(args ...string) error {
	if len(args) > 0 {
		return fmt.Errorf("ls doesn't take any argument")
	}
	files, err := ioutil.ReadDir(store)
	if err != nil {
		return fmt.Errorf("listing files: %s", err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil
}

func (cli *Cli) help(args ...string) error {
	if len(args) > 0 {
		return fmt.Errorf("doesn't take any argument yet")
	}

	fmt.Println("Money")
	fmt.Println("\nAvailable commands:")
	fmt.Println("exit")
	for cmd, _ := range cli.commands {
		fmt.Println(cmd)
	}

	return nil
}

func (cli *Cli) unhandledCommand(command string, args []string) {
	// FIXME: look up similar commands
	fmt.Printf("Command %q doesn't exist\n", command)
}

func parse(instructionline string) (string, []string) {
	// FIXME: this might need to be fancier (support quotes and stuff).
	fields := strings.Fields(instructionline)
	return fields[0], fields[1:]
}

func NewCli(cryptor *Cryptor) *Cli {
	return &Cli{
		cryptor: cryptor,
	}
}
