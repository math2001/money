package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/math2001/money/api"
)

func main() {
	fmt.Println("Welcome to Money!")
	fmt.Println("=================")
	fmt.Println()

	api, err := api.NewAPI("data")
	if err != nil {
		log.Fatalf("Creating api: %s", err)
	}

	r := mux.NewRouter()
	html := r.Methods(http.MethodGet).Subrouter()

	html.PathPrefix("/css").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./pwa/css"))))
	html.HandleFunc("/js/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Service-Worker-Allowed", "/")
		http.ServeFile(w, r, "/js/sw.js")
	})
	html.PathPrefix("/js").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./pwa/js"))))

	api.BindTo(r.PathPrefix("/api").Subrouter())

	html.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving html to %q %q", r.Method, r.URL)
		http.ServeFile(w, r, "pwa/index.html")
	})

	http.Handle("/", r)
	log.Printf("Ready. Listening on :9999")
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
