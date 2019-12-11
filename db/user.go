package db

import (
	"fmt"
	"os"
)

// Maybe rename to db.go?

// UserDB is the folder containing all the user's data.
// Everything in their is encrypted using his password
type User struct {
	// root is the user's own folder (see api.go)
	root    string
	email   string
	cryptor *Cryptor
}

func (u *User) Save(filename string, plaintext []byte) error {
	panic("not implemented")
}

func (u *User) Load(filename string) ([]byte, error) {
	panic("not implemented")
}

func (u *User) Login(password []byte) error {
	panic("not implemented")
}

func (u *User) SignUp(password []byte) error {
	// FIXME: safety check: make sure that privroot doesn't already exists.
	if err := os.MkdirAll(u.root, 0644); err != nil {
		return fmt.Errorf("signing up, creating user folder: %s", err)
	}

	panic("not implemented")
}

func NewUser(root, email string) *User {
	return &User{
		root:  root,
		email: email,
	}
}
