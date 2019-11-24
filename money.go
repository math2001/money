package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	fmt.Println("Welcome to Cryptor")

	// FIXME: how can we make sure that the user can easily check that *this*
	// program is asking for the password, and some other random thing?

	var cryptor *Cryptor
	// cryptor = login()

	if err := os.MkdirAll("store", os.ModePerm); err != nil {
		log.Fatalf("makdir store: %s", err)
	}

	NewApp(cryptor).Start()
}

func login() *Cryptor {
	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("reading password from stdin: %s", err)
	}

	cryptor, err := NewCryptor(password)
	if err != nil {
		log.Fatalf("creating cryptor: %s", err)
	}
	return cryptor
}
