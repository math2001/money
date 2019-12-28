package server

import (
	"errors"
	"image/png"
	"log"
	"net/http"

	"github.com/math2001/money/api"
)

func (s *Server) addManualPayment(r *http.Request) *resp {
	user, err := s.getCurrentUser(r)
	if errors.Is(err, ErrNoCurrentUser) {
		log.Printf("add payments: no current user")
		return &resp{
			code: http.StatusNotAcceptable,
			body: kv{
				"kind":    "require log in",
				"msg":     "please authenticate first",
				"details": "authentication cookie found, but user forgotten",
			},
		}
	} else if err != nil {
		log.Printf("[err] loading session: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			body: kv{
				"kind": "not acceptable",
				"msg":  "couldn't load session from cookie",
			},
		}
	}

	err = s.api.AddPayment(user, []byte(r.PostFormValue("payment")))

	if _, ok := err.(api.ErrInvalidPayment); ok {
		log.Printf("invalid payment: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			body: kv{
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
			body: kv{
				"kind": "internal error",
				"msg":  "adding payment failed",
			},
		}
	}
	return &resp{
		code: http.StatusOK,
		body: kv{
			"kind": "success",
			"goto": "/", // FIXME: where should it go
		},
	}
}

func (s *Server) listPayments(r *http.Request) *resp {

	user, err := s.getCurrentUser(r)
	if errors.Is(err, ErrNoCurrentUser) {
		log.Printf("add payments: no current user")
		return &resp{
			code: http.StatusNotAcceptable,
			body: kv{
				"kind":    "require log in",
				"msg":     "please authenticate first",
				"details": "authentication cookie found, but user forgotten",
			},
		}
	} else if err != nil {
		log.Printf("[err] loading session: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			body: kv{
				"kind": "not acceptable",
				"msg":  "couldn't load session from cookie",
			},
		}
	}

	payments, err := s.api.ListPayments(user)
	if err != nil {
		log.Printf("[err] listing payments: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			body: kv{
				"kind": "internal error",
			},
		}
	}

	return &resp{
		code: http.StatusOK,
		body: kv{
			"kind":     "success",
			"payments": payments,
		},
	}
}

func (s *Server) scan(r *http.Request) *resp {

	user, err := s.getCurrentUser(r)
	if errors.Is(err, ErrNoCurrentUser) {
		log.Printf("scan: no current user")
		return &resp{
			code: http.StatusNotAcceptable,
			body: kv{
				"kind":    "require log in",
				"msg":     "please authenticate first",
				"details": "authentication cookie found, but user forgotten",
			},
		}
	} else if err != nil {
		log.Printf("[err] scan: loading session: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			body: kv{
				"kind": "not acceptable",
				"msg":  "couldn't load session from cookie",
			},
		}
	}

	file, header, err := r.FormFile("img")
	if err != nil {
		log.Printf("[err] scan: loading file from post requets: %s", err)
		return &resp{
			code: http.StatusBadRequest,
			body: kv{
				"kind": "bad request",
				"msg":  "error occurred when uploading file",
			},
		}
	}
	defer file.Close()
	// FIXME: SECURITY check file size before uploading
	// FIXME: check that it's the right image format

	img, err := png.Decode(file)
	if err != nil {
		log.Printf("[err] png.Decoding file")
		return &resp{
			code: http.StatusExpectationFailed,
			body: kv{
				"kind":  "expectation failed",
				"msg":   "couldn't PNG decode file",
				"FIXME": "support multiple file format",
			},
		}
	}

	payment, err := s.api.Scan(user, header, img)
	if err != nil {
		log.Printf("[err] listing payments: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			body: kv{
				"kind": "internal error",
				"id":   "scan job failed",
			},
		}
	}

	return &resp{
		code: http.StatusOK,
		body: kv{
			"kind":    "success",
			"payment": payment,
		},
	}
}
