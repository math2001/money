package api

import (
	"errors"
	"fmt"
	"log"
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

	user, err := api.Login(r.PostFormValue("email"), r.PostFormValue("password"))
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

	user, err := api.SignUp(email, password)
	if errors.Is(err, ErrEmailAlreadyUsed) {
		respond(w, r, http.StatusNotAcceptable, "email already used")
		return
	} else if err != nil {
		log.Printf("[err] signing up: %s", err)
		respond(w, r, http.StatusInternalServerError, "error", "msg", "failed to sign up user")
		return
	}

	fmt.Println("got user", user)
	// write http cookie

	// FIXME: check session for where the user is coming from, and redirect him
	// there (don't forget to remove that session item)
	respond(w, r, http.StatusOK, "goto", "goto", "/")

}
