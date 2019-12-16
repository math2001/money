package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/math2001/money/api"
)

func main() {
	fmt.Println("Welcome to Money!")
	fmt.Println("=================")
	fmt.Println()

	handler := getHandler("data", os.Stdout)

	log.Printf("Ready. Listening on :9999")
	if err := http.ListenAndServe(":9999", handler); err != nil {
		log.Fatal(err)
	}
}

func getHandler(dataroot string, logs io.Writer) http.Handler {
	api, err := api.NewAPI(dataroot)
	if err != nil {
		log.Fatalf("Creating api: %s", err)
	}

	log.SetOutput(logs)

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
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusPermanentRedirect)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// FIXME: implement warning system
			log.Printf("serving %q GET request with html", r.URL.Path)
		}
		http.ServeFile(w, r, "pwa/index.html")
	})

	r.Use(logger)

	return r
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%q %q", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
