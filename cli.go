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

// the private folder name in which we store the keys, salts and password hash
// FIXME: allow the user to change this
const privroot = "priv"

// Start the CLI application
func (cli *Cli) Start() {

	cli.commands = map[string]func(...string) error{
		"help": cli.help,
		"ls":   cli.ls,
		"load": cli.load,
		"save": cli.save,
	}

	var err error
	cli.km = NewKeysManager(privroot)
	if err != nil {
		log.Fatalf("create keys manager: %s", err)
	}

	if cli.km.HasSignedUp() {
		fmt.Println("Log in")
		if err := cli.login(); err != nil {
			log.Fatalf("logging in: %s", err)
		}
	} else {
		fmt.Println("We will create your account, please sign up")
		if err := cli.signup(); err != nil {
			log.Fatalf("signing up: %s", err)
		}
	}

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

func (cli *Cli) login() error {
	password, err := cli.askpassword("Enter password")
	if err != nil {
		return err
	}

	err = cli.km.Login(password)
	if errors.Is(err, ErrNoSaltsFile) {
		return fmt.Errorf("Your private folder is corrupted (%w)", err)
	}
	if err != nil {
		return fmt.Errorf("login in keys manager: %s", err)
	}

	keys, err := cli.km.LoadKeys()
	if errors.Is(err, ErrNoKeysfile) {
		return fmt.Errorf("Your private folder is corrupted (%w)", err)
	}

	if err != nil {
		return fmt.Errorf("loading keys: %s", err)
	}

	cli.cryptor, err = NewCryptor(keys.MAC, keys.Encryption)
	if err != nil {
		return fmt.Errorf("creating cryptor: %s", err)
	}
	return nil
}

func (cli *Cli) signup(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("signup doesn't take any arguments")
	}

	password, err := cli.askpassword("Enter password")
	if err != nil {
		return err
	}
	confirm, err := cli.askpassword("Confirm password")
	if err != nil {
		return err
	}

	if !bytes.Equal(password, confirm) {
		// FIXME: should this be an ErrPasswordsDontMatch?
		return fmt.Errorf("passwords don't match, please try again")
	}

	if err := cli.km.SignUp(password); err != nil {
		return fmt.Errorf("signing up: %s", err)
	}

	keys, err := cli.km.LoadKeys()
	if err != nil {
		return fmt.Errorf("loading keys: %s", err)
	}

	cli.cryptor, err = NewCryptor(keys.MAC, keys.Encryption)
	if err != nil {
		return fmt.Errorf("creating cryptor: %s", err)
	}

	fmt.Println("successfully created your account!")
	return nil
}

func (cli *Cli) askpassword(message string) ([]byte, error) {
	fmt.Printf("%s: ", message)
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("reading password from stdin: %s", err)
	}
	return password, nil
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
	fmt.Print(string(content))
	if content[len(content)-1] != '\n' {
		// funny character to indicate that we manually added a line return
		// idea and character stolen from the fish shell
		fmt.Println("\xe2\x8f\x8e\x20")
	}
	return nil
}

func (cli *Cli) save(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("takes one argument, the filename")
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

func (cli *Cli) unhandledCommand(command string, args []string) {
	// FIXME: look up similar commands
	fmt.Printf("Command %q doesn't exist\n", command)
}

func parse(instructionline string) (string, []string) {
	// FIXME: this might need to be fancier (support quotes and stuff).
	fields := strings.Fields(instructionline)
	return fields[0], fields[1:]
}
