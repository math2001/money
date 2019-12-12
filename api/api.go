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
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/math2001/money/keysmanager"
	"github.com/math2001/money/sessions"
)

const saltsize = 32

const (
	// saltcookie is used to sign the cookie's payload
	saltcookie = iota

	// saltpasswod is used encrypt the passwords within the database
	saltpassword

	// saltsession is used to sign the cookies
	// FIXME: this should be a key, and not a salt, but that'll mean that the
	// application will need it's own keymanager, which leads to 2
	// possibilities:
	// 1. ask for a user input (password) everytime the api starts (= a pain)
	// 2. store the keysmanager password in a file, which is equivalent to
	// storing the keys as a salt (from a technical point of view, because this
	// shouldn't be shared at all, unlike other salts who just need to be
	// unique)
	// 2nd solution is much better because we get a sense of which file are
	// important to keep secret.
	saltsession

	// the number of salts I need
	_nsalts
)

type API struct {
	dataroot  string
	userslist string
	usersdir  string
	// this salt is used to hash the passwords in the database
	sm       *keysmanager.SM
	sessions *sessions.S
}

func NewAPI(dataroot string) (*API, error) {
	log.Printf("API dataroot: %q", dataroot)

	saltfile := filepath.Join(dataroot, "appsalts")
	sm := keysmanager.NewSaltsManager(_nsalts, saltfile, saltsize)

	// the datafoot folder doesn't exists
	if _, err := os.Stat(dataroot); os.IsNotExist(err) {
		log.Println("initiating fresh api...")
		if err := os.Mkdir(dataroot, 0700); err != nil {
			return nil, fmt.Errorf("mkdir %q: %s", dataroot, err)
		}
		if err := sm.GenerateNew(); err != nil {
			return nil, fmt.Errorf("generating new salts: %s", err)
		}
	} else {
		log.Printf("Resuming from filesystem...")
		if err := sm.Load(); err != nil {
			return nil, fmt.Errorf("loading salt: %s", err)
		}
	}

	s, err := sessions.NewS(&sessions.Config{
		Key: sm.Get(saltsession),
	})

	if err != nil {
		return nil, fmt.Errorf("creating sessions.S: %s", err)
	}

	return &API{
		sessions:  s,
		sm:        sm,
		dataroot:  dataroot,
		userslist: filepath.Join(dataroot, "users.list"),
		usersdir:  filepath.Join(dataroot, "users"),
	}, nil
}

// Serve starts a http server under /api/
func (api *API) BindTo(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if r.URL.Path != "/api/" {
		// 	respond(w, r, http.StatusNotFound, "endpoint undefined")
		// 	return
		// }
		respond(w, r, http.StatusOK, "FIXME: list all the possible endpoints")
	})

	post := r.Methods(http.MethodPost).Subrouter()

	post.HandleFunc("/login", api.loginHandler)
	post.HandleFunc("/signup", api.signupHandler)
}
