package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/math2001/money/api"
)

func (s *Server) login(r *http.Request) *resp {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	user, err := s.api.Login(email, password)

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
			},
		}
	}

	encryptedPassword, err := s.cryptor.Encrypt([]byte(password))
	if err != nil {
		log.Printf("[err] encrypting password: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
			},
		}
	}

	return &resp{
		code: http.StatusOK,
		session: &Session{
			ID:       user.ID,
			Email:    user.Email,
			Password: encryptedPassword,
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

	r.ParseMultipartForm(1 << 20) // 1 MB of memory

	// FIXME: validate email and password
	// email: regex and stdlib? or just strings.contains
	// password: min length, evaluate strength?

	emails, eok := r.PostForm["email"]
	passwords, pok := r.PostForm["password"]
	confirms, cok := r.PostForm["confirm"]
	if !eok || !pok || !cok {
		log.Printf("!! warning !! missing fields %t %t %t", eok, pok, cok)
		return &resp{
			code: http.StatusBadRequest,
			msg: kv{
				"kind": "bad request",
				"msg":  "missing fields",
				"help": []string{
					"This is the fault of the developers, not yours. Please report this issue :-)",
				},
			},
		}
	}

	if len(emails) != 1 || len(passwords) != 1 || len(confirms) != 1 {
		log.Printf("!! warning !! duplicate fields %d %d %d", len(emails), len(passwords), len(confirms))
		return &resp{
			code: http.StatusBadRequest,
			msg: kv{
				"kind": "bad request",
				"msg":  "duplicate fields",
				"help": []string{
					"This is the fault of the developers, not yours. Please report this issue :-)",
				},
			},
		}
	}

	email, password, confirm := emails[0], passwords[0], confirms[0]

	err := checkmail.ValidateFormat(email)
	if errors.Is(err, checkmail.ErrBadFormat) {
		return &resp{
			code: http.StatusExpectationFailed,
			msg: kv{
				"kind": "invalid input",
				"msg":  fmt.Sprintf("Invalid email: %q didn't match our required format", email),
			},
		}
	} else if err != nil {
		log.Printf("!! warning !! unknown error from checkmail: %s", err)
		return &resp{
			code: http.StatusExpectationFailed,
			msg: kv{
				"kind": "invalid input",
				"msg":  "invalid email: unkown error",
			},
		}
	}

	if password != confirm {
		return &resp{
			code: http.StatusExpectationFailed,
			msg: kv{
				"kind": "password dismatch",
				"msg":  "Your passwords don't match. Please try again.",
			},
		}
	}

	if len(password) < 8 {
		return &resp{
			code: http.StatusExpectationFailed,
			msg: kv{
				"kind": "password too short",
				"msg":  "Your password is too short. It should be at least 8 characters.",
			},
		}
	}

	// FIXME: check passwords against common bank of password

	user, err := s.api.SignUp(email, password)
	if errors.Is(err, api.ErrEmailAlreadyUsed) {
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "email already used",
				"msg":  "This emails has already been used. Did you forget your password? If yes, you should contact us",
				"FIXME": []string{
					"implement password reset via email",
				},
			},
		}
	} else if err != nil {
		log.Printf("[err] signing up: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
			},
		}
	}

	encryptedPassword, err := s.cryptor.Encrypt([]byte(password))
	if err != nil {
		log.Printf("[err] encrypting password: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
			},
		}
	}

	// FIXME: check session for where the user is coming from, and redirect him
	// there (don't forget to remove that session item)

	return &resp{
		code: http.StatusOK,
		session: &Session{
			ID:       user.ID,
			Email:    user.Email,
			Password: encryptedPassword,
		},
		msg: kv{
			"kind":  "success",
			"goto":  "/",
			"email": email,
		},
	}
}

// logout always deletes the session, regardless of the error the it encouters
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

	user, err := s.getCurrentUser(r)
	if errors.Is(err, ErrNoCurrentUser) {
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "not acceptable",
				"msg":  "no user is logged in",
			},
			session: &NilSession,
		}
	} else if err != nil {
		log.Printf("[err] loading session: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "not acceptable",
				"msg":  "couldn't load session from cookie",
			},
			session: &NilSession,
		}
	}

	pwaEmail := r.PostFormValue("email")
	if pwaEmail != user.Email {
		log.Printf("!! warning !! the pwa's email doesn't match the session's email")
		log.Printf("\npwa email: %q\nsession email: %q", pwaEmail, user.Email)
		// FIXME: should we let the user know about this?
	}

	err = s.api.Logout(user)
	if err != nil {
		log.Printf("logout handler: api.Logout: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
			},
			session: &NilSession,
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
