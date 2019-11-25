package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var commands = map[string]func(...string){
	"generate-hex-keys": generateHexKey,
}

func main() {
	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}
	cmd, ok := commands[os.Args[1]]
	if !ok {
		fmt.Printf("Command %q not found\n", os.Args[1])
	}
	cmd(os.Args[2:]...)
}

func usage() {
	fmt.Println("$ tools <cmd> <args>")
	fmt.Println("Commands: ")
	for cmd := range commands {
		fmt.Println(cmd)
	}
}

// FIXME: print in a go code way, so that we can just copy paste
func generateHexKey(args ...string) {
	if len(args) != 1 {
		fmt.Println("generate-hex-keys nkeys")

	}

	nkeys, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Parsing args: %s", err)
	}

	const keysize = 32

	for i := 0; i < nkeys; i++ {
		key := make([]byte, keysize)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			log.Fatalf("generating key: %s", err)
		}
		fmt.Println(hex.EncodeToString(key))
	}
}
