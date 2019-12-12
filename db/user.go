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
	ID      int
	Email   string
	root    string
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
	if err := os.Mkdir(u.root, 0700); err != nil {
		return fmt.Errorf("signing up, creating user folder: %s", err)
	}
	return nil
}

func NewUser(id int, email, root string) *User {
	return &User{
		root:  root,
		Email: email,
		ID:    id,
	}
}
