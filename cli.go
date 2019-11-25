package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// Cli is the cli application which provides an interface to the Cryptor
type Cli struct {
	cryptor  *Cryptor
	commands map[string]func(...string) error
	reader   *bufio.Reader
	km       *KeysManager
}

// Start the CLI application
func (cli *Cli) Start() {
	cli.commands = map[string]func(...string) error{
		"login": cli.login,
		"help":  cli.help,
		"ls":    cli.ls,
		"load":  cli.load,
		"save":  cli.save,
	}

	cli.km = NewKeysManager()

	cli.reader = bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		instructionline, err := cli.reader.ReadString('\n')
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

func (cli *Cli) login(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("login doesn't take any argument")
	}
	// FIXME: check km.HasKeysfile before asking for a password
	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("reading password from stdin: %s", err)
	}

	if err := cli.km.Login(password); err != nil {
		return fmt.Errorf("login in keys manager: %s", err)
	}

	keys, err := cli.km.LoadKeys()
	if errors.Is(err, ErrNoKeysfile) {
		fmt.Println("No keysfile found.")
		if !cli.confirm("Would you like to generate some new keys?") {
			return fmt.Errorf("Abort.")
		}
		cli.generatenewkeys()
		fmt.Println("You now have to login")
		return cli.login()
	}
	if err != nil {
		return fmt.Errorf("loading keys from password: %s", err)
	}

	cli.cryptor, err = NewCryptor(keys.MAC, keys.Encryption)
	if err != nil {
		return fmt.Errorf("creating cryptor: %s", err)
	}
	return nil
}

func (cli *Cli) load(args ...string) error {
	if len(args) != 1 {
		// FIXME: display usage for this command
		return fmt.Errorf("load takes one argument, the filename")
	}
	if err := cli.IsLoggedIn(); err != nil {
		return err
	}
	path := filepath.Join(store, args[0])
	content, err := cli.cryptor.Load(path)
	if err != nil {
		return fmt.Errorf("loading %q: %s", path, err)
	}
	fmt.Print(string(content))
	if content[len(content)-1] != '\n' {
		fmt.Println("\xe2\x8f\x8e\x20")
	}
	return nil
}

func (cli *Cli) save(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("takes one argument, the filename")
	}

	if err := cli.IsLoggedIn(); err != nil {
		return err
	}

	path := filepath.Join(store, args[0])
	// file exists and user doesn't want to overwrite it
	if _, err := os.Stat(path); err == nil && !cli.confirm("Overwrite existing file?") {
		fmt.Println("Abort")
		return nil
	}

	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading from stdin: %s", err)
	}
	if err := cli.cryptor.Save(path, content); err != nil {
		return fmt.Errorf("saving to file: %s", err)
	}
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
	for cmd := range cli.commands {
		fmt.Println(cmd)
	}

	return nil
}

func (cli *Cli) confirm(question string) bool {
	for {
		fmt.Printf("%s (y/n) ", question)
		ans, err := cli.reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			log.Fatalf("reading line: %s", err)
		}
		if ans == "y\n" {
			return true
		} else if ans == "n\n" {
			return false
		}
		// otherwise we keep asking
	}
}

func (cli *Cli) generatenewkeys() error {
	fmt.Print("Enter new password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("reading password from stdin: %s", err)
	}

	fmt.Print("Confirm new password: ")
	confirm, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("reading confirm from stdin: %s", err)
	}

	if !bytes.Equal(confirm, password) {
		return fmt.Errorf("Passwords don't match")
	}

	if err := cli.km.GenerateNewKeys(password); err != nil {
		// FIXME: instead of deleting the keys file, rename to back up file,
		// like keys#.priv where # is the backup number
		if errors.Is(err, ErrKeysfileExists) {
			if !cli.confirm("overwrite keysfile? existing keys will be lost forever") {
				fmt.Println("Abort")
				return nil
			}
			if err := cli.km.RemoveKeysfile(); err != nil {
				return err
			}
			if err := cli.km.GenerateNewKeys(password); err != nil {
				return fmt.Errorf("generating new keys: %s", err)
			}
		} else {
			return fmt.Errorf("generate new keys: %s", err)
		}
	}
	fmt.Println("new keys generated successfully")
	return nil
}

func (cli *Cli) unhandledCommand(command string, args []string) {
	// FIXME: look up similar commands
	fmt.Printf("Command %q doesn't exist\n", command)
}

func (cli *Cli) IsLoggedIn() error {
	if cli.cryptor == nil {
		return fmt.Errorf("You need to login first\n    > login")
	}
	return nil
}

func parse(instructionline string) (string, []string) {
	// FIXME: this might need to be fancier (support quotes and stuff).
	fields := strings.Fields(instructionline)
	return fields[0], fields[1:]
}
