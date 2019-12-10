package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./pwa/css"))))

	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/js/sw.js" {
			w.Header().Set("Service-Worker-Allowed", "/")
		}
		http.ServeFile(w, r, filepath.Join("./pwa", r.URL.Path))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pwa/index.html")
	})

	log.Printf("listening on :9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}
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

var ErrOddParts = errors.New("cannot generate map from odd number of parts")
var ErrReservedKey = errors.New("reserved key")

func responde(w http.ResponseWriter, kind string, parts ...interface{}) error {
	if len(parts)%2 == 1 {
		return fmt.Errorf("# parts: %d (%w)", len(parts), ErrOddParts)
	}

	obj := make(map[string]interface{}, len(parts)/2)
	obj["kind"] = kind
	for i, part := range parts {
		if i%2 == 1 {
			continue
		}
		if key, ok := part.(string); ok {
			if key == "kind" {
				return fmt.Errorf("'kind' (%w)", ErrReservedKey)
			}
			obj[key] = parts[i+1]
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(obj); err != nil {
		return fmt.Errorf("writing json obj: %s", err)
	}

	return nil
}
