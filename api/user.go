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
	"golang.org/x/crypto/scrypt"
)

var ErrEmailAlreadyUsed = errors.New("email already used")

var ErrWrongIdentifiers = errors.New("wrong identifiers")

var ErrNoCurrentUser = errors.New("no current user")

type user struct {
	Email    string
	Password []byte
	ID       int
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
	log.Printf("Sign up %q", email)

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

	// FIXME: the key size (32) should be a constant
	hashedpassword, err := scrypt.Key([]byte(password), api.sm.Get(saltpassword), 32768, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("signing up, hashing password: %s", err)
	}

	// add entry in users file
	userid := len(users) + 1
	users = append(users, user{
		Email:    email,
		Password: hashedpassword,
		ID:       userid,
	})

	f.Close()
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

	u := db.NewUser(userid, email, filepath.Join(api.usersdir, strconv.Itoa(userid)))
	if err := u.SignUp([]byte(password)); err != nil {
		return nil, fmt.Errorf("signing up db.User: %s", err)
	}

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

	u := db.NewUser(match.ID, match.Email, filepath.Join(api.usersdir, strconv.Itoa(match.ID)))
	if err := u.Login([]byte(password)); err != nil {
		return nil, fmt.Errorf("logging in: %s", err)
	}

	log.Printf("user is now logged in: %v", u)

	return u, nil
}

// Logout has nothing to do to log out someone from the api's point of view
// so we just at least check that the current user is valid
func (api *API) Logout(id int, email string) error {
	_, err := api.getCurrentUser(id, email)
	if err != nil {
		return err
	}
	return nil
}

// getCurrentUser returns the current user. error can be ErrNoCurrentUser,
// ErrInvalidUser
func (api *API) getCurrentUser(id int, email string) (*db.User, error) {
	panic("api.getCurrentUser not implemented")
	// u, ok := api.loggedusers[id]
	// if !ok {
	// 	return nil
	// }
	// if u.ID != id {
	// 	log.Printf("!! warning !! loggedusers is broken: actual id: %d, key: %d", u.ID, id)
	// 	// this is a major issue, hence require logging in
	// 	delete(api.loggedusers, id)
	// 	return nil
	// }
	// if u.Email != email {
	// 	log.Printf("!! warning !! current user %d %q doesn't have expected email %q", u.ID, u.Email, email)
	// }

	// return u
}
