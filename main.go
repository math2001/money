package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/math2001/money/api"
)

func main() {
	fmt.Println("Welcome to Money!")
	fmt.Println("=================")
	fmt.Println()

	r := startAt("data")

	http.Handle("/", r)
	log.Printf("Ready. Listening on :9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}
}

func startAt(dataroot string) *mux.Router {
	api, err := api.NewAPI(dataroot)
	if err != nil {
		log.Fatalf("Creating api: %s", err)
	}
	_ = api

	r := mux.NewRouter().StrictSlash(true)
	api.BindTo(r.PathPrefix("/api").Subrouter())
	html := r.Methods(http.MethodGet).Subrouter()

	html.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./pwa/css"))))
	html.HandleFunc("/js/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Service-Worker-Allowed", "/")
		http.ServeFile(w, r, "/js/sw.js")
	})
	html.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./pwa/js"))))

	html.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// FIXME: implement warning system
			log.Printf("serving %q GET request with html", r.URL.Path)
		}
		log.Printf("Serving html to %q %q", r.Method, r.URL)
		http.ServeFile(w, r, "pwa/index.html")
	})

	return r
}
