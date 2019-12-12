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
	"io/ioutil"
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

// Session is the content of the session cookie
type Session struct {
	ID    int
	Email string
}

func NewAPI(dataroot string) (*API, error) {
	log.Printf("API dataroot: %q", dataroot)

	api := &API{
		dataroot:  dataroot,
		userslist: filepath.Join(dataroot, "users.list"),
		usersdir:  filepath.Join(dataroot, "users"),
	}

	saltfile := filepath.Join(dataroot, "appsalts")
	api.sm = keysmanager.NewSaltsManager(_nsalts, saltfile, saltsize)

	// the datafoot folder doesn't exists, start from scratch
	if _, err := os.Stat(dataroot); os.IsNotExist(err) {
		log.Println("initiating fresh api...")
		if err := os.Mkdir(api.dataroot, 0700); err != nil {
			return nil, fmt.Errorf("mkdir %q: %s", api.dataroot, err)
		}
		if err := os.Mkdir(api.usersdir, 0700); err != nil {
			return nil, fmt.Errorf("mkdir %q: %s", api.usersdir, err)
		}
		if err := api.sm.GenerateNew(); err != nil {
			return nil, fmt.Errorf("generating new salts: %s", err)
		}
		if err := ioutil.WriteFile(api.userslist, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("writing {} file")
		}

	} else {
		log.Printf("Resuming from filesystem...")
		if err := api.sm.Load(); err != nil {
			return nil, fmt.Errorf("loading salt: %s", err)
		}
	}

	var err error
	api.sessions, err = sessions.NewS(&sessions.Config{
		Key: api.sm.Get(saltsession),
	})

	if err != nil {
		return nil, fmt.Errorf("creating sessions.S: %s", err)
	}
	return api, nil
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
