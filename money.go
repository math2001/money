package main

import (
	"fmt"
	"log"
	"os"
)

const store = "store"

func main() {
	fmt.Println("Welcome to Money")

	// FIXME: how can we make sure that the user can easily check that *this*
	// program is asking for the password, and some other random thing?

	if err := os.MkdirAll(store, os.ModePerm); err != nil {
		log.Fatalf("makdir store: %s", err)
	}

	cli := &Cli{}
	cli.Start()
}
