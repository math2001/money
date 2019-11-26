package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

type SaltsManager struct {
	saltsize int
	saltdir  string
	Salts    *Salts
}

var ErrNoSaltFile = errors.New("no salt file")
var ErrAlreadyLoaded = errors.New("salts already loaded")

// salts containts the *decrypted* salts
// we use a struct instead of a map to have error checking at compile time, but
// it comes at the cost of having hacks in load/save function to make sure that
// we save/load every field (and reflecting field name is really hacky).
// It's better to have a simple and safe interface with a few hacks in the
// implementation rather than a simple implementation with a risky interface.
// But maybe it isn't worth it (ie we should use a map)
type Salts struct {
	cipher, password []byte
}

func (s *Salts) Cipher() []byte   { return s.cipher }
func (s *Salts) Password() []byte { return s.password }

var ErrNilSalt = errors.New("nil salt")

func (sm *SaltsManager) GenerateNewSalts() error {
	generateNewSalt := func() ([]byte, error) {
		salt := make([]byte, sm.saltsize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, fmt.Errorf("generating salt: %s", err)
		}
		return salt, nil
	}

	cipher, err := generateNewSalt()
	if err != nil {
		return err
	}
	password, err := generateNewSalt()
	if err != nil {
		return err
	}

	sm.Salts = &Salts{cipher: cipher, password: password}

	// keep it split up because we will store the salts in a single file

	saveSalt := func(name string, salt []byte) error {
		if salt == nil {
			return fmt.Errorf("saving %q: %w", name, ErrNilSalt)
		}

		path := filepath.Join(sm.saltdir, name)

		// FIXME: use copy and wrap into a hex.NewEncoder. It won't noticeably
		// improve perfs, but it's cleaner
		if err := ioutil.WriteFile(path, []byte(hex.EncodeToString(salt)), 0644); err != nil {
			// CHECKME: should we give the location of the file (ie. path)
			return fmt.Errorf("writing salt %q: %s", name, err)
		}
		return nil
	}

	// a hack to make sure we don't forget to save the salt
	if reflect.TypeOf(sm.Salts).NumField() != 2 {
		// if you get this error, you need to make sure that you save the new
		// key that has been just added: just below, there should be an other
		// call saveSalt("<new key>"). Then, you can increment the number in
		// the condition above. This number above should be the number saveSalt
		// calls
		panic("[internal] a new salt has been added, but the save function wasn't updated")
	}

	if err := saveSalt("cipher", sm.Salts.Cipher()); err != nil {
		return err
	}
	if err := saveSalt("password", sm.Salts.Password()); err != nil {
		return err
	}
	return nil
}

func (sm *SaltsManager) LoadSalts() error {
	if sm.Salts != nil {
		return ErrAlreadyLoaded
	}

	getSalt := func(name string) ([]byte, error) {
		path := filepath.Join(sm.saltdir, name)

		hexsalt, err := ioutil.ReadFile(path)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%q %w", name, ErrNoSaltFile)
		}
		if err != nil {
			return nil, fmt.Errorf("reading salt file: %s", err)
		}
		salt, err := hex.DecodeString(string(hexsalt))
		if err != nil {
			return nil, fmt.Errorf("decoding salt: %s", err)
		}
		return salt, nil
	}

	// a hack to make sure we don't forget to save the salt. See Save method
	if reflect.TypeOf(sm.Salts).NumField() != 2 {
		panic("[internal] a new salt has been added, but the save function wasn't updated")
	}
	cipher, err := getSalt("cipher")
	if err != nil {
		return err
	}
	password, err := getSalt("password")
	if err != nil {
		return err
	}

	sm.Salts = &Salts{
		cipher:   cipher,
		password: password,
	}

	return nil
}

// RefreshSalts replaces the current salts with new ones, but keeps everything
// working. I'm too lazy right now to think about whether or not it's even
// possible
func (sm *SaltsManager) RefreshSalts() {
	panic("not implemented")
}

// NewSaltsManager returns a salt manager with some sane defaults
func NewSaltsManager(privroot string) *SaltsManager {
	return &SaltsManager{
		saltsize: 16,
		saltdir:  filepath.Join(privroot, "salts"),
	}
}
