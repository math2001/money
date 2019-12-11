package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/math2001/money/api"
)

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

	mux := &http.ServeMux{}
	mux.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./pwa/css"))))
	mux.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		// the service worker needs a special header because it is served from
		// ./js/ (hence it's max scope is ./js/), but I need it's scope to be /
		if r.URL.Path == "/js/sw.js" {
			w.Header().Set("Service-Worker-Allowed", "/")
		}

		http.ServeFile(w, r, filepath.Join("./pwa", r.URL.Path))
	})

	api, err := api.NewAPI("data")
	if err != nil {
		log.Fatalf("creating api: %s", err)
	}
	api.Serve(mux)

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

	log.Fatalf("not implemented")

	// FIXME: how can we make sure that the user can easily check that *this*
	// program is asking for the password, and some other random thing?

	// cli := &Cli{}
	// cli.Start()
}
