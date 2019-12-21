package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/math2001/money/keysmanager"
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
	Usersdir  string
	// this salt is used to hash the passwords in the database
	sm *keysmanager.SM
}

func NewAPI(dataroot string) *API {
	log.Printf("API dataroot: %q", dataroot)

	api := &API{
		dataroot:  dataroot,
		userslist: filepath.Join(dataroot, "users.list"),
		Usersdir:  filepath.Join(dataroot, "users"),
	}

	saltfile := filepath.Join(dataroot, "apisalts")
	api.sm = keysmanager.NewSaltsManager(_nsalts, saltfile, saltsize)

	return api
}

// Initialize creates all the required file (should only be run if they don't
// already exist)
func (api *API) Initialize() error {
	if err := os.Mkdir(api.Usersdir, 0700); err != nil {
		return fmt.Errorf("mkdir %q: %s", api.Usersdir, err)
	}
	if err := api.sm.GenerateNew(); err != nil {
		return fmt.Errorf("generating new salts: %s", err)
	}
	if err := ioutil.WriteFile(api.userslist, []byte("[]"), 0644); err != nil {
		return fmt.Errorf("writing [] to file %s", err)
	}
	return nil
}

// Resume loads stuff from files to be ready to serve
func (api *API) Resume() error {
	if err := api.sm.Load(); err != nil {
		return fmt.Errorf("loading salt: %s", err)
	}
	return nil
}
