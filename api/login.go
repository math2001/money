package api

import (
	"log"
	"net/http"
)

type loginInfos struct {
	id       int
	username int
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if err := respond(w, http.StatusMethodNotAllowed, "method not allowed", "method", r.Method); err != nil {
			log.Printf("%v writing method not allowed: %s", r, err)
		}
		return
	}

}
