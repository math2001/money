package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/math2001/money/db"
)

type ErrorReportStore struct {
	*db.Store
}

type reportFrom string

const (
	ReportFromUser   reportFrom = "report from user"
	ReportFromServer reportFrom = "report from server"
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

var ErrUnauthorized = errors.New("unauthorized")

// probably should be general to api if not more
var ErrNotFound = errors.New("not found")

// Report adds the report to the reports :^)
func (api *API) Report(report *Report) error {
	// TODO: send an email to admin when there is a user report
	if api.errorreports == nil {
		log.Printf("%v", report)
		return errors.New("error reporter hasn't been initiated. Dropped report to logs")
	}
	return api.errorreports.Add(report)
}

// ReportsList lists all the reports.
func (api *API) ReportsList(user *db.User) ([]string, error) {
	if !user.Admin {
		return nil, ErrUnauthorized
	}
	filenames, err := api.errorreports.List()
	if err != nil {
		return nil, fmt.Errorf("getting reports list: %s", err)
	}
	// FIXME: maybe just load the type/description/date
	return filenames, nil
}

// GetReport retrieves a report from the filename.
func (api *API) GetReport(user *db.User, filename string) (*Report, error) {
	if !user.Admin {
		return nil, ErrUnauthorized
	}
	return api.errorreports.Get(filename)
}

// Add adds the report (JSON encoded) to a new file (filename = timestamp.json)
func (ers *ErrorReportStore) Add(report *Report) error {
	log.Printf("new error report %q", report.Kind)
	filename := fmt.Sprintf("%d.json", time.Now().UnixNano())
	for i := 0; ers.Exists(filename); i++ {
		filename = fmt.Sprintf("%d.json", time.Now().UnixNano())
		if i > 100 {
			// FIXME: tag internal error
			return errors.New("couldn't generate report name (already exist)")
		}
	}

	report.Date = time.Now()

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

func (ers *ErrorReportStore) Get(filename string) (*Report, error) {
	content, err := ers.Load(filename)
	var patherr *os.PathError
	if ok := errors.As(err, &patherr); ok && os.IsNotExist(patherr) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var report *Report
	err = json.Unmarshal(content, &report)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling report %q: %s", filename, err)
	}
	return report, nil
}

func NewErrorReportStore(root string) *ErrorReportStore {
	return &ErrorReportStore{
		db.NewStore(root),
	}
}
