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
		respond(w, r, http.StatusMethodNotAllowed, "method not allowed", "method", r.Method)
		return
	}

	user, err := api.Login(r.PostFormValue("email"), r.PostFormValue("password"))
	if errors.Is(err, ErrWrongIdentifiers) {
		respond(w, r, http.StatusOK, "wrong identification")
		return
	} else if err != nil {
		log.Printf("[err] loging in: %s", err)
		respond(w, r, http.StatusOK, "error", "msg", "internal error while logging in")
		return
	}

	api.sessions.Save(w, &Session{
		ID:    user.ID,
		Email: user.Email,
	})

	respond(w, r, http.StatusOK, "success", "goto", "/")
}

func (api *API) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond(w, r, http.StatusMethodNotAllowed, "method not allowed", "method", r.Method)
		return
	}

	// FIXME: validate email and password
	// email: regex and stdlib? or just strings.contains
	// password: min length, evaluate strength?

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

	if err := api.sessions.Save(w, &Session{
		ID:    user.ID,
		Email: user.Email,
	}); err != nil {
		log.Printf("[err] saving session: %s", err)
		respond(w, r, http.StatusInternalServerError, "error", "msg", "failed to sign up user")
		return
	}

	// FIXME: check session for where the user is coming from, and redirect him
	// there (don't forget to remove that session item)
	respond(w, r, http.StatusOK, "success", "goto", "/", "email", email)
}
