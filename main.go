package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/math2001/money/server"
)

func main() {
	fmt.Println("Welcome to Money!")
	fmt.Println("=================")
	fmt.Println()

	// FIXME: ask for password from stdin
	password := []byte("")

	handler, err := server.New("data", password)
	if err != nil {
		log.Fatalf("creating server: %s", err)
	}

	log.Printf("Ready. Listening on :9999")
	if err := http.ListenAndServe(":9999", handler); err != nil {
		log.Fatal(err)
	}
}
