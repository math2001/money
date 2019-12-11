package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/math2001/money/db"
	"golang.org/x/crypto/scrypt"
)

var ErrEmailAlreadyUsed = errors.New("email already used")

var ErrWrongIdentifiers = errors.New("wrong identifiers")

// SignUp creates a new user
//
// FIXME: this function can change the state of the application but still
// return an error. It needs to clean up after itself if that happens, because
// otherwise, we are left with a corrupted state
func (api *API) SignUp(email, password []byte) (*db.User, error) {
	// check taken entry in users file
	f, err := os.Open(api.userslist)
	if err != nil {
		// TODO: could send an email to email (apologise)
		return nil, fmt.Errorf("signing up, opening users list %q: %s", api.userslist, err)
	}
	decoder := json.NewDecoder(f)

	// FIXME: this is extrememly inefficient. It reads all the user data into
	// memory just to compare emails and possibly add one entry
	type user struct {
		email    []byte
		password []byte
		id       int
	}

	var users []user
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("signing up, parsing users list: %q: %s", api.userslist, err)
	}

	// check if the email has already been used
	for _, user := range users {
		if bytes.Equal(user.email, email) {
			return nil, ErrEmailAlreadyUsed
		}
	}

	// FIXME: the key size (32) should be a constant
	hashed, err := scrypt.Key(password, api.sm.Get(saltpassword), 32768, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("signing up, hashing password: %s", err)
	}

	// add entry in users file

	userid := len(users) + 1
	users = append(users, struct {
		email    []byte
		password []byte
		id       int
	}{
		email:    email,
		password: hashed,
		id:       userid,
	})

	// WHAT? how does that work? the file was open in read mode...
	encoder := json.NewEncoder(f)
	if err := encoder.Encode(users); err != nil {
		// FIXME: the users file is now corrupted. Try to rewrite the old
		// version
		return nil, fmt.Errorf("signing up, saving user to database: %s", err)
	}

	u := db.NewUser(filepath.Join(api.usersdir, strconv.Itoa(userid)), string(email))
	u.SignUp(password)
	return u, nil
}

func (api *API) Login(email, password []byte) (*db.User, error) {

	// FIXME: that's a lot of duplicate logic from sign up...

	f, err := os.Open(api.userslist)
	if err != nil {
		// TODO: could send an email to email (apologise)
		return nil, fmt.Errorf("signing up, opening users list %q: %s", api.userslist, err)
	}
	decoder := json.NewDecoder(f)

	// FIXME: this is extrememly inefficient. It reads all the user data into
	// memory just to compare emails and password pairs...
	type user struct {
		email    []byte
		password []byte
		id       int
	}

	var users []user
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("signing up, parsing users list: %q: %s", api.userslist, err)
	}

	hashed, err := scrypt.Key(password, api.sm.Get(saltpassword), 32768, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("signing up, hashing password: %s", err)
	}

	var match user
	// check if the email has already been used
	for _, user := range users {
		if bytes.Equal(user.email, email) && bytes.Equal(user.password, hashed) {
			match = user
			break
		}
	}

	// ie. no match
	if match.id == 0 {
		return nil, ErrWrongIdentifiers
	}

	u := db.NewUser(filepath.Join(api.usersdir, strconv.Itoa(match.id)), string(email))
	u.Login(password)
	return u, nil
}
