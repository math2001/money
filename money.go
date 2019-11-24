package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// FIXME: how can we make sure that the user can easily check that *this*
	// program is asking for the password, and some other random thing?
	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("reading password from stdin: %s", err)
	}

	cryptor, err := NewCryptor(password)
	if err != nil {
		log.Fatalf("creating cryptor: %s", err)
	}

	if err := os.MkdirAll("store", os.ModePerm); err != nil {
		log.Fatalf("makdir store: %s", err)
	}

	input := []byte("aaaaa bbbb cccc dddd")

	if err := cryptor.Save("store/test1", input); err != nil {
		log.Fatalf("saving to store/test1: %s", err)
	}

	output, err := cryptor.Load("store/test1")
	if err != nil {
		log.Fatalf("loading from store/test1: %s", err)
	}

	fmt.Printf("%q\n%q\nEqual: %t", input, output, bytes.Equal(input, output))
}
