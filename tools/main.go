package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("need one argument: <action>")
	}

	actions := map[string]func(){
		"parsesession": ParseSession,
		"ps":           ParseSession,
	}
	action, ok := actions[os.Args[1]]
	if !ok {
		log.Fatalf("unknown action: %q", os.Args[1])
	}
	action()
}