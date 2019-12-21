package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/math2001/money/keysmanager"
)

// Maybe rename to db.go?

// UserDB is the folder containing all the user's data.
// Everything in their is encrypted using his password
type User struct {
	// root is the user's own folder (see api.go)
	ID          int
	Email       string
	root        string
	cryptor     *Cryptor
	keysmanager *keysmanager.KeysManager
}

func (u *User) Save(filename string, plaintext []byte) error {
	path := JoinRootPath(u.root, filename)
	return u.cryptor.Save(path, plaintext)
}

func (u *User) Load(filename string) ([]byte, error) {
	path := JoinRootPath(u.root, filename)
	return u.cryptor.Load(path)
}

// Login can return keysmanager.ErrWrongPassword, keysmanager.ErrPrivCorrupted,
// ErrAlreadyLoaded (internal) or err
func (u *User) Login(password []byte) error {
	err := u.keysmanager.Login(password)
	if errors.Is(err, keysmanager.ErrWrongPassword) || errors.Is(err, keysmanager.ErrPrivCorrupted) {
		return err
	} else if errors.Is(err, keysmanager.ErrAlreadyLoaded) {
		// FIXME: tag internal
		return err
	} else if err != nil {
		return fmt.Errorf("db.login keysm.Login: %s", err)
	}

	keys, err := u.keysmanager.LoadKeys()
	if errors.Is(err, keysmanager.ErrPrivCorrupted) {
		return err
	} else if err != nil {
		return fmt.Errorf("db.login keysm.LoadKeys: %s", err)
	}

	c, err := NewCryptor(keys.Encryption, keys.MAC)
	if err != nil {
		return fmt.Errorf("db.login newcryptor: %s", err)
	}
	u.cryptor = c

	return nil
}

func (u *User) SignUp(password []byte) error {
	if err := os.Mkdir(u.root, 0700); err != nil {
		return fmt.Errorf("signing up, creating user folder: %s", err)
	}
	// FIXME: initiate the cryptor

	err := u.keysmanager.SignUp(password)
	if errors.Is(err, keysmanager.ErrAlreadyLoaded) {
		return err
	} else if err != nil {
		// FIXME: tag internal
		return err
	}

	keys, err := u.keysmanager.LoadKeys()
	if errors.Is(err, keysmanager.ErrPrivCorrupted) {
		return err
	} else if err != nil {
		return fmt.Errorf("db.signup keysm.LoadKeys: %s", err)
	}

	c, err := NewCryptor(keys.Encryption, keys.MAC)
	if err != nil {
		return fmt.Errorf("db.signup newcryptor: %s", err)
	}
	u.cryptor = c

	return nil
}

func NewUser(id int, email, root string) *User {
	return &User{
		root:        root,
		Email:       email,
		ID:          id,
		keysmanager: keysmanager.NewKeysManager(filepath.Join(root, "secrets")),
	}
}
