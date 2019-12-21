package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/math2001/money/api"
)

func (s *Server) login(r *http.Request) *resp {
	email := r.PostFormValue("email")
	user, err := s.api.Login(email, r.PostFormValue("password"))

	if errors.Is(err, api.ErrWrongIdentifiers) {
		return &resp{
			code: http.StatusOK, // FIXME: better error code?
			msg: kv{
				"kind": "wrong identifiers",
			},
		}
	} else if err != nil {
		log.Printf("[err] logging in: %s", err)
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
			"kind":  "success",
			"goto":  "/",
			"email": user.Email,
		},
	}
}

func (s *Server) signup(r *http.Request) *resp {
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

	user, err := s.api.SignUp(email, password)
	if errors.Is(err, api.ErrEmailAlreadyUsed) {
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

func (s *Server) logout(r *http.Request) *resp {
	r.ParseMultipartForm(1024) // 1 KB

	if _, ok := r.PostForm["email"]; !ok {
		return &resp{
			code: http.StatusBadRequest,
			msg: kv{
				"kind": "bad request",
				"msg":  "expected 'email' field",
			},
		}
	}

	// now, we compare that the email that was stored by the pwa is indeed the
	// email of the current user this is just a safety check because doing
	// offline is going to be hard if they don't match, the user (wrong email
	// and correct) will just be logged out of both things. But we know that
	// there was a problem somewhere...

	// remember that sessionEmail can be trusted because it's payload is signed

	session, err := s.getSession(r)
	if err != nil {
		log.Printf("[err] loading session: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "not acceptable",
				"msg":  "couldn't load session from cookie",
			},
		}
	}

	pwaEmail := r.PostFormValue("email")
	if pwaEmail != session.Email {
		log.Printf("!! warning !! the pwa's email doesn't match the session's email")
		log.Printf("\npwa email: %q\nsession email: %q", pwaEmail, session.Email)
		// FIXME: should we let the user know about this?
	}

	err = s.api.Logout(session.ID, session.Email)
	if errors.Is(err, api.ErrNoCurrentUser) {
		return &resp{
			code: http.StatusExpectationFailed,
			msg: kv{
				"kind": "error",
				"id":   "no user",
				"msg":  "no user is currently logged in",
			},
			session: &NilSession, // remove the session cookie
		}
	} else if err != nil {
		log.Printf("logout handler: api.Logout: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
				"msg":  "logging user out",
			},
		}
	}

	return &resp{
		code: http.StatusOK,
		msg: kv{
			"kind": "success",
			"goto": "/",
		},
		session: &NilSession, // remove the session cookie
	}
}
