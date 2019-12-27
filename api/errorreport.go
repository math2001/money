package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/math2001/money/db"
)

type ErrorReportStore struct {
	*db.Store
}

type reportFrom int

const (
	ReportFromUser reportFrom = iota
	ReportFromServer
)

type Report struct {
	Kind           string
	Description    string
	From           reportFrom
	Date           time.Time
	User           *db.User
	ErrGettingUser error
	Err            error
	Request        *http.Request `json:"-"`
	Details        map[string]interface{}

	// RequestURI should be ignored by the user. Just give the whole request
	// object
	RequestURI string
}

// Add adds the report (JSON encoded) to a new file (filename = timestamp.json)
func (ers *ErrorReportStore) add(report *Report) error {
	log.Printf("new error report %q", report.Kind)
	filename := fmt.Sprintf("%d.json", time.Now().UnixNano())
	for i := 0; ers.Exists(filename); i++ {
		filename = fmt.Sprintf("%d.json", time.Now().UnixNano())
		if i > 100 {
			// FIXME: tag internal error
			return errors.New("couldn't generate report name (already exist)")
		}
	}

	report.RequestURI = report.Request.RequestURI
	report.Request = nil
	content, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshaling report: %s", err)
	}

	if err := ers.Save(filename, content); err != nil {
		return fmt.Errorf("saving report: %s", err)
	}
	return nil
}

func NewErrorReportStore(root string) *ErrorReportStore {
	return &ErrorReportStore{
		db.NewStore(root),
	}
}

// Report adds the report to the reports :^)
func (api *API) Report(report *Report) error {
	// TODO: send an email to admin when there is a user report
	if api.errorreports == nil {
		log.Printf("%v", report)
		return errors.New("error reporter hasn't been initiated. Dropped report to logs")
	}
	return api.errorreports.add(report)
}
