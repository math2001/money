package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/math2001/money/api"
)

func (s *Server) addManualPayment(r *http.Request) *resp {
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

	err = s.api.AddPayment(session.ID, session.Email, []byte(r.PostFormValue("payment")))

	if errors.Is(err, api.ErrNoCurrentUser) {
		log.Printf("add payments: no current user")
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind":    "require log in",
				"msg":     "please authenticate first",
				"details": "authentication cookie found, but user forgotten",
			},
		}
	} else if _, ok := err.(api.ErrInvalidPayment); ok {
		log.Printf("invalid payment: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "error",
				"id":   "invalid payment",
				"msg":  err.Error(),
			},
		}
	}
	if err != nil {
		log.Printf("add payments: api.addpayment: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
				"msg":  "adding payment failed",
			},
		}
	}
	return &resp{
		code: http.StatusOK,
		msg: kv{
			"kind": "success",
			"goto": "/", // FIXME: where should it go
		},
	}
}
