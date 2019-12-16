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

func (api *API) loginHandler(r *http.Request) *resp {
	user, err := api.Login(r.PostFormValue("email"), r.PostFormValue("password"))

	if errors.Is(err, ErrWrongIdentifiers) {
		return &resp{
			code: http.StatusOK, // FIXME: better error code?
			msg: kv{
				"kind": "wrong identifiers",
			},
		}
	} else if err != nil {
		log.Printf("[err] loging in: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
				"msg":  "logging in",
			},
		}
	}

	return &resp{
		code: http.StatusOK,
		session: &Session{
			ID:    user.ID,
			Email: user.Email,
		},
		msg: kv{
			"kind": "success",
			"goto": "/",
		},
	}
}

func (api *API) signupHandler(r *http.Request) *resp {
	if r.Method != http.MethodPost {
		return &resp{
			code: http.StatusMethodNotAllowed,
			msg: kv{
				"kind":   "method not allowed",
				"method": r.Method,
			},
		}
	}

	// FIXME: validate email and password
	// email: regex and stdlib? or just strings.contains
	// password: min length, evaluate strength?

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	confirm := r.PostFormValue("confirm")
	if password != confirm {
		return &resp{
			code: http.StatusExpectationFailed,
			msg: kv{
				"kind": "password dismatch",
			},
		}
	}

	user, err := api.SignUp(email, password)
	if errors.Is(err, ErrEmailAlreadyUsed) {
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "email already used",
			},
		}
	} else if err != nil {
		log.Printf("[err] signing up: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
				"msg":  "failed to sign up user",
			},
		}
	}

	// FIXME: check session for where the user is coming from, and redirect him
	// there (don't forget to remove that session item)

	return &resp{
		code: http.StatusOK,
		session: &Session{
			ID:    user.ID,
			Email: user.Email,
		},
		msg: kv{
			"kind":  "success",
			"goto":  "/",
			"email": email,
		},
	}
}
