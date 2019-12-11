package api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var ErrEmailAlreadyUsed = errors.New("email already used")

type API struct {
	dataroot  string
	userslist string
	usersdir  string
	salt      []byte
}

const saltsize = 32

func NewAPI(dataroot string) (*API, error) {
	if err := os.MkdirAll(dataroot, 0700); err != nil {
		return nil, fmt.Errorf("mkdir dataroot: %s", err)
	}

	saltfile := filepath.Join(dataroot, "salts")

	salt, err := ioutil.ReadFile(saltfile)
	if os.IsNotExist(err) {
		salt = make([]byte, saltsize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, fmt.Errorf("generating salt: %s", err)
		}
		if err := ioutil.WriteFile(saltfile, salt, 0600); err != nil {
			return nil, fmt.Errorf("writing salt to file: %s", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("opening saltfile %q: %s", saltfile, err)
	}

	return &API{
		salt:      salt,
		dataroot:  dataroot,
		userslist: filepath.Join(dataroot, "users.list"),
		usersdir:  filepath.Join(dataroot, "users"),
	}, nil
}

// Serve starts a http server under /api/
func (api *API) BindTo(mux *http.ServeMux) {
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, "FIXME: list all the possible actions")
	})

	mux.HandleFunc("/api/login", loginHandler)
}
