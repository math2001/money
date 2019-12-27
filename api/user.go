package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/math2001/money/db"
)

var ErrEmailAlreadyUsed = errors.New("email already used")

var ErrWrongIdentifiers = errors.New("wrong identifiers")

// the fields have to be exported for json to be able to access them when
// encoding
type user struct {
	Email    string
	Password []byte
	ID       int
	Admin    bool
}

// SignUp creates a new user
//
// FIXME: this function can change the state of the application but still
// return an error. It needs to clean up after itself if that happens, because
// otherwise, we are left with a corrupted state
//
// FIXME: this is extrememly inefficient. It reads all the user data into
// memory just to compare emails and possibly add one entry
func (api *API) SignUp(email, password string) (*db.User, error) {
	// check taken entry in users file
	f, err := os.Open(api.userslist)
	if err != nil {
		// TODO: could send an email to email (apologise)
		return nil, fmt.Errorf("signing up, opening users list %q: %s", api.userslist, err)
	}
	decoder := json.NewDecoder(f)

	var users []user
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("signing up, parsing users list: %q: %s", api.userslist, err)
	}

	// check if the email has already been used
	for _, user := range users {
		if user.Email == email {
			log.Printf("Return email already used")
			return nil, ErrEmailAlreadyUsed
		}
	}

	hashedpassword := scryptKey([]byte(password), api.sm.Get(saltpassword))

	// add entry in users file
	userid := len(users) + 1
	newuser := user{
		Email:    email,
		Password: hashedpassword,
		ID:       userid,
		// the first user to sign up is automatically admin
		Admin: len(users) == 0,
	}
	users = append(users, newuser)

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("closing users.list file: %s", err)
	}

	f, err = os.Create(api.userslist)
	if err != nil {
		return nil, fmt.Errorf("signing up, recreating users list: %s", err)
	}
	encoder := json.NewEncoder(f)
	if err := encoder.Encode(users); err != nil {
		// FIXME: the users file is now corrupted. Try to rewrite the old
		// version
		return nil, fmt.Errorf("signing up, saving user to database: %s", err)
	}

	u := db.NewUser(newuser.ID, newuser.Email, newuser.Admin, filepath.Join(api.Usersdir, strconv.Itoa(userid)))
	if err := u.SignUp([]byte(password)); err != nil {
		return nil, fmt.Errorf("signing up db.User: %s", err)
	}

	log.Printf("Signed up %q", newuser.Email)

	return u, nil
}

// Login adds the user to loggedusers
func (api *API) Login(email, password string) (*db.User, error) {

	// FIXME: that's a lot of duplicate logic from sign up...

	f, err := os.Open(api.userslist)
	if err != nil {
		// TODO: could send an email to email (apologise)
		return nil, fmt.Errorf("signing up, opening users list %q: %s", api.userslist, err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)

	// FIXME: this is extrememly inefficient. It reads all the user data into
	// memory just to compare emails and password pairs...
	var users []user
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("signing up, parsing users list: %q: %s", api.userslist, err)
	}

	hashedpassword := scryptKey([]byte(password), api.sm.Get(saltpassword))

	var match user
	// check if the email has already been used
	for _, user := range users {
		if user.Email == email && bytes.Equal(user.Password, hashedpassword) {
			match = user
			break
		}
	}

	// ie. no match
	if match.ID == 0 {
		return nil, ErrWrongIdentifiers
	}

	u := db.NewUser(match.ID, match.Email, match.Admin, filepath.Join(api.Usersdir, strconv.Itoa(match.ID)))
	if err := u.Login([]byte(password)); err != nil {
		return nil, fmt.Errorf("logging in: %s", err)
	}

	log.Printf("Logged in %q", u.Email)

	return u, nil
}

// Logout has nothing to do to log out someone from the api's point of view
// so we just at least check that the current user is valid
func (api *API) Logout(user *db.User) error {
	return nil
}
