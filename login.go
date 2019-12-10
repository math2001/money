package main

import (
	"log"
	"net/http"
)

type loginInfos struct {
	id       int
	username int
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := responde(w, "method not allowed"); err != nil {
			log.Printf("%v writing method not allowed: %s", r, err)
		}
	}

}
