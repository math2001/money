package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const store = "store"

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Invalid number of arguments: \n   $ money <mode>")
	}
	if os.Args[0] == "cli" {
		CLIMode()
	} else if os.Args[1] == "server" {
		ServerMode()
	} else {
		log.Fatalf("Invalid mode: %q. Only support cli or server", os.Args[1])
	}
}

func ServerMode() {
	fmt.Println("Welcome to Money! [server mode]")
	fmt.Println("===============================")
	fmt.Println()

	http.Handle("/", http.FileServer(http.Dir("./pwa/")))
	log.Printf("listening on :9999")
	http.ListenAndServe(":9999", nil)

}

func CLIMode() {
	fmt.Println("Welcome to Money! [cli mode]")
	fmt.Println("============================")
	fmt.Println()

	// FIXME: how can we make sure that the user can easily check that *this*
	// program is asking for the password, and some other random thing?

	if err := os.MkdirAll(store, os.ModePerm); err != nil {
		log.Fatalf("makdir store: %s", err)
	}

	cli := &Cli{}
	cli.Start()
}
