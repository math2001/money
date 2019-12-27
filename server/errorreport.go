package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/math2001/money/api"
)

func (s *Server) reportsNew(r *http.Request) *resp {
	errorreport := []byte(r.PostFormValue("report"))

	user, erruser := s.getCurrentUser(r)
	if erruser != nil && !errors.Is(erruser, ErrNoCurrentUser) {
		log.Printf("!! warning !! getting current user for error report: %s", erruser)
		// keep going, we will just have no user for that error report
	}

	var report *api.Report
	var details map[string]interface{}
	if err := json.Unmarshal(errorreport, &details); err != nil {
		// don't save the report as plain text, because it's probably something
		// malicious happening here. If it's during dev, then just print out
		// the errorreport, but in prod, this is probably an attack
		log.Printf("!! warning !! Potential attack: unmarshalling error report: %s", err)

		report = &api.Report{
			Kind:           "failed user report",
			Description:    "Couldn't JSON decode the user report. This is a potential attack",
			From:           api.ReportFromUser,
			User:           user,
			ErrGettingUser: erruser,
			Err:            err,
		}
	} else {
		report = &api.Report{
			Kind:           "user report",
			From:           api.ReportFromUser,
			User:           user,
			ErrGettingUser: erruser,
			Details:        details,
		}
	}

	if err := s.api.Report(report); err != nil {
		log.Printf("[err] api.ErrorReport: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
			},
		}
	}
	return nil
}

func (s *Server) reportsList(r *http.Request) *resp {
	user, err := s.getCurrentUser(r)
	if errors.Is(err, ErrNoCurrentUser) {
		log.Printf("listing reports: no current user")
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "require log in",
				"msg":  "please authenticate first",
			},
		}
	} else if err != nil {
		log.Printf("[err] loading session: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "not acceptable",
				"msg":  "couldn't load session from cookie",
			},
		}
	}

	reports, err := s.api.ReportsList(user)

	if errors.Is(err, api.ErrUnauthorized) {
		return &resp{
			code: http.StatusUnauthorized,
			msg: kv{
				"kind": "unauthorized",
				"msg":  "you are not authorized to view this page",
			},
		}
	} else if err != nil {
		// FIXME: put err in resp, and parse it out in the handler
		log.Printf("[err] internal error: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
			},
		}
	}

	return &resp{
		code: http.StatusOK,
		msg: kv{
			"kind":    "success",
			"reports": reports,
		},
	}
}

func (s *Server) reportsGet(r *http.Request) *resp {
	user, err := s.getCurrentUser(r)
	if errors.Is(err, ErrNoCurrentUser) {
		log.Printf("getting reports: no current user")
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "unauthorized",
				"msg":  "please authenticate first",
			},
		}
	} else if err != nil {
		log.Printf("[err] loading session: %s", err)
		return &resp{
			code: http.StatusNotAcceptable,
			msg: kv{
				"kind": "not acceptable",
				"msg":  "couldn't load session from cookie",
			},
		}
	}

	filenames, _ := r.URL.Query()["filename"]
	if len(filenames) != 1 {
		return &resp{
			code: http.StatusBadRequest,
			msg: kv{
				"kind": "bad request",
				"msg":  "requires one 'filename' field",
			},
		}
	}

	report, err := s.api.GetReport(user, filenames[0])

	if errors.Is(err, api.ErrUnauthorized) {
		return &resp{
			code: http.StatusUnauthorized,
			msg: kv{
				"kind": "unauthorized",
				"msg":  "you are not authorized to view this page",
			},
		}
	} else if errors.Is(err, api.ErrNotFound) {
		return &resp{
			code: http.StatusNotFound,
			msg: kv{
				"kind": "not found",
				"msg":  fmt.Sprintf("report %q not found", filenames[0]),
			},
		}
	} else if err != nil {
		// FIXME: put err in resp, and parse it out in the handler
		log.Printf("[err] internal error: %s", err)
		return &resp{
			code: http.StatusInternalServerError,
			msg: kv{
				"kind": "internal error",
			},
		}
	}

	return &resp{
		code: http.StatusOK,
		msg: kv{
			"kind":   "success",
			"report": report,
		},
	}
}
