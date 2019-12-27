package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/math2001/money/api"
)

func (s *Server) reporterror(r *http.Request) *resp {
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

	report.Date = time.Now()

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
