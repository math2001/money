package api

import (
	"errors"
	"net/http"
)

type loginInfos struct {
	id       int
	username int
}

func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond(w, r, http.StatusMethodNotAllowed, "method not allowed", "method", r.Method)
		return
	}

	user, err := api.Login([]byte(r.PostFormValue("email")), []byte(r.PostFormValue("password")))
	if errors.Is(err, ErrWrongIdentifiers) {
		respond(w, r, http.StatusOK, "wrong identification")
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
