package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/crypto/scrypt"
)

var ErrEmailAlreadyUsed = errors.New("email already used")

// UserDB is the folder containing all the user's data.
// Everything in their is encrypted using his password
type UserDB struct {
	privroot string
}

type App struct {
	salt      []byte
	root      string
	userslist string
	usersdir  string
}

func NewApp(root string) (*App, error) {
	// FIXME: the app should have it's own password (which would be required on
	// start up), just like a regular user. The salt file and userslist could
	// then be encrypted,

	salt, err := ioutil.ReadFile(filepath.Join(root, "salt"))
	if os.IsNotExist(err) {
		// generate new salt
		// FIXME: this part should be exposed so that we can use deterministic
		// salts for testing
	}
	if err != nil {
		return nil, fmt.Errorf("reading app salt: %s", err)
	}

	return &App{
		salt:      salt,
		root:      root,
		userslist: filepath.Join(root, "users.list"),
		usersdir:  filepath.Join(root, "users"),
	}, nil
}

// SignUp creates a new user
// FIXME: this function can change the state of the application but still
// return an error. It needs to clean up after itself if that happens, because
// otherwise, we are left with a corrupted state
func (app *App) SignUp(email, password []byte) (*UserDB, error) {
	// check taken entry in users file
	f, err := os.Open(app.userslist)
	if err != nil {
		// TODO: could send an email to that guy...
		return nil, fmt.Errorf("signing up, opening users list %q: %s", app.userslist, err)
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
		return nil, fmt.Errorf("signing up, parsing users list: %q: %s", app.userslist, err)
	}

	// check if the email has already been used
	for _, user := range users {
		if bytes.Equal(user.email, email) {
			return nil, ErrEmailAlreadyUsed
		}
	}

	hashed, err := scrypt.Key(password, app.salt, 32768, 8, 1, 32)
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

	// FIXME: safety check: make sure that privroot doesn't already exists.
	privroot := filepath.Join(app.usersdir, strconv.Itoa(userid))
	if err := os.MkdirAll(privroot, 0644); err != nil {
		return nil, fmt.Errorf("signing up, creating user folder: %s", err)
	}
	return &UserDB{
		privroot: privroot,
	}, nil
}

// func Login(password []byte) (*UserDB, error) {
// 	// check entry matches (return nil on sucess)
// }

// func (u *UserDB) Save(filename string, plaintext []byte) error {

// }

// func (u *UserDB) Load(filename string) ([]byte, error) {

// }
