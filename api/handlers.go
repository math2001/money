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
	respond(w, r, http.StatusOK, "success")
}

func (api *API) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond(w, r, http.StatusMethodNotAllowed, "method not allowed", "method", r.Method)
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	confirm := r.PostFormValue("confirm")
	if password != confirm {
		respond(w, r, http.StatusExpectationFailed, "password dismatch")
		return
	}

	user, err := api.SignUp([]byte(email), []byte(password))
	if errors.Is(err, ErrEmailAlreadyUsed) {
		respond(w, r, http.StatusNotAcceptable, "email already used")
	}

	_ = user
	// write http cookie

	// FIXME: check session for where the user is coming from, and redirect him
	// there (don't forget to remove that session item)
	respond(w, r, http.StatusOK, "goto", "goto", "/")

}
