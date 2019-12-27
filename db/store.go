package db

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/math2001/money/keysmanager"
)

// Store is a folder containing all the store's data.
// Everything in their is encrypted using the store's password
type Store struct {
	root        string
	cryptor     *Cryptor
	keysmanager *keysmanager.KeysManager
}

var ErrAuthenticateFirst = errors.New("authentication required")

func (s *Store) Save(filename string, plaintext []byte) error {
	path := JoinRootPath(s.root, filename)
	if s.cryptor == nil {
		return ErrAuthenticateFirst
	}
	return s.cryptor.Save(path, plaintext)
}

func (s *Store) Load(filename string) ([]byte, error) {
	path := JoinRootPath(s.root, filename)
	if s.cryptor == nil {
		return nil, ErrAuthenticateFirst
	}
	return s.cryptor.Load(path)
}

// List returns the names of all the files
func (s *Store) List() ([]string, error) {
	files, err := ioutil.ReadDir(s.root)
	if err != nil {
		return nil, err
	}
	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = file.Name()
	}
	return filenames, nil
}

func (s *Store) Exists(filename string) bool {
	_, err := os.Stat(JoinRootPath(s.root, filename))
	return err == nil || os.IsExist(err)
}

// Login can return keysmanager.ErrWrongPassword, keysmanager.ErrPrivCorrupted,
// ErrAlreadyLoaded (internal) or err
func (s *Store) Login(password []byte) error {
	err := s.keysmanager.Login(password)
	if errors.Is(err, keysmanager.ErrWrongPassword) || errors.Is(err, keysmanager.ErrPrivCorrupted) {
		return err
	} else if errors.Is(err, keysmanager.ErrAlreadyLoaded) {
		// FIXME: tag internal
		return err
	} else if err != nil {
		return fmt.Errorf("db.login keysm.Login: %s", err)
	}

	keys, err := s.keysmanager.LoadKeys()
	if errors.Is(err, keysmanager.ErrPrivCorrupted) {
		return err
	} else if err != nil {
		return fmt.Errorf("db.login keysm.LoadKeys: %s", err)
	}

	c, err := NewCryptor(keys.Encryption, keys.MAC)
	if err != nil {
		return fmt.Errorf("db.login newcryptor: %s", err)
	}
	s.cryptor = c

	return nil
}

func (s *Store) SignUp(password []byte) error {
	if err := os.MkdirAll(s.root, 0700); err != nil {
		return fmt.Errorf("signing up, creating store folder: %s", err)
	}

	err := s.keysmanager.SignUp(password)
	if errors.Is(err, keysmanager.ErrAlreadyLoaded) {
		return err
	} else if err != nil {
		// FIXME: tag internal
		return err
	}

	keys, err := s.keysmanager.LoadKeys()
	if errors.Is(err, keysmanager.ErrPrivCorrupted) {
		return err
	} else if err != nil {
		return fmt.Errorf("db.signup keysm.LoadKeys: %s", err)
	}

	c, err := NewCryptor(keys.Encryption, keys.MAC)
	if err != nil {
		return fmt.Errorf("db.signup newcryptor: %s", err)
	}
	s.cryptor = c

	return nil
}

func (s Store) String() string {
	return fmt.Sprintf("Store{root: %q}", s.root)
}

func NewStore(root string) *Store {
	return &Store{
		root:        filepath.Join(root, "data"),
		keysmanager: keysmanager.NewKeysManager(filepath.Join(root, "secrets")),
	}
}
