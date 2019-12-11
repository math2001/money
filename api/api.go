/*
Data folder structure:

    appsalts
    users/
        {id1}/
            kmdata/                # this is managed by keys manager
                keys
                passwordhashfile   # in the user's case, this file is redundant
                salts
            recordfile1            # these are encrypted with kmdata just above
            recordfile2
            ...
        {id2}/
            kmdata/
                keys
                passwordhashfile
                salts
            recordfile1
            recordfile2
            ...


    The appsalts contains two salts:

    saltcookie                 # the salt used to generate the keys used to
                               # sign the cookies
    saltpassword               # the salt used to encrypt the passwords
                               # within the database
*/
package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/math2001/money/keysmanager"
)

type API struct {
	dataroot  string
	userslist string
	usersdir  string
	// this salt is used to hash the passwords in the database
	sm *keysmanager.SM
}

const saltsize = 32

const (
	// saltcookie is used to sign the cookie's payload
	saltcookie = iota

	// saltpasswod is used encrypt the passwords within the database
	saltpassword
)

func NewAPI(dataroot string) (*API, error) {
	saltfile := filepath.Join(dataroot, "appsalts")
	sm := keysmanager.NewSaltsManager(2, saltfile, saltsize)
	sm.Load()
	return &API{
		sm:        sm,
		dataroot:  dataroot,
		userslist: filepath.Join(dataroot, "users.list"),
		usersdir:  filepath.Join(dataroot, "users"),
	}, nil
}

// IsUninitiated checks if the API has already been running before (ie. there are
// files it can resume serving from)
func (api *API) IsUninitiated() bool {
	if _, err := os.Stat(api.dataroot); os.IsNotExist(err) {
		return true
	}
	return false
}

func (api *API) Initiate() error {
	// FIXME: clean up after yourself if you return early (error)
	if err := os.Mkdir(api.dataroot, 0600); err != nil {
		return fmt.Errorf("mkdir %q: %s", api.dataroot, err)
	}
	if err := api.sm.GenerateNew(); err != nil {
		return err
	}
	return nil
}

// Serve starts a http server under /api/
func (api *API) BindTo(mux *http.ServeMux) {
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, "FIXME: list all the possible actions")
	})

	mux.HandleFunc("/api/login", api.loginHandler)
}
