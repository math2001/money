package api

import (
	"errors"
	"log"
	"net/http"
)

type loginInfos struct {
	id       int
	username int
}

func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if err := respond(w, http.StatusMethodNotAllowed, "method not allowed", "method", r.Method); err != nil {
			log.Printf("%v writing method not allowed: %s", r, err)
		}
		return
	}

	user, err := api.Login([]byte(r.PostFormValue("email")), []byte(r.PostFormValue("password")))
	if errors.Is(err, ErrWrongIdentifiers) {
		if err := respond(w, http.StatusOK, "wrong identification"); err != nil {
			log.Printf("%v writing wrong id: %s", r, err)
		}
		return
	}

	_ = user

	// write HTTP cookie
	api.sm.Get(saltcookie)

	// write success
	if err := respond(w, http.StatusOK, "success"); err != nil {
		log.Printf("%v writing success: %s", r, err)
	}

}
