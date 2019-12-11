package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/math2001/money/api"
)

func main() {
	fmt.Println("Welcome to Money! [server mode]")
	fmt.Println("===============================")
	fmt.Println()

	api, err := api.NewAPI("data")
	if err != nil {
		log.Fatalf("creating api: %s", err)
	}

	if api.IsUninitiated() {
		fmt.Println("api is fresh (no data folder). Initiating now...")
		api.Initiate()
		fmt.Println("FIXME: give user a choice")
	}

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

	api.BindTo(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pwa/index.html")
	})

	log.Printf("listening on :9999")
	if err := http.ListenAndServe(":9999", mux); err != nil {
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
