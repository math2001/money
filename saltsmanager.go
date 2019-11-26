package main

import (
	"bufio"
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
	saltsize  int
	saltsfile string
	Salts     *Salts
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

// FIXME: use a map internally, and then expose a static, compile-time checked
// interface

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

	// a hack to make sure we don't forget to save the salt
	if reflect.TypeOf(sm.Salts).NumField() != 2 {
		// if you get this error, you need to make sure that you save the new
		// key that has been just added: just below, there should be an other
		// call saveSalt("<new key>"). Then, you can increment the number in
		// the condition above. This number above should be the number saveSalt
		// calls
		panic("[internal] a new salt has been added, but the save function wasn't updated")
	}

	salts := fmt.Sprintf("%s\n%s\n", sm.Salts.Cipher(), sm.Salts.Password())
	if err := ioutil.WriteFile(sm.saltsfile, []byte(salts), 0644); err != nil {
		return fmt.Errorf("writing salts to file system: %s", err)
	}

	return nil
}

func (sm *SaltsManager) LoadSalts() error {
	if sm.Salts != nil {
		return ErrAlreadyLoaded
	}

	f, err := os.Open(sm.saltsfile)
	if err != nil {
		return fmt.Errorf("opening saltsfile: %s", err)
	}
	reader := bufio.NewReader(f)

	hexcipher, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading cipher salt: %s", err)
	}

	hexpassword, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading password salt: %s", err)
	}

	cipher, err := hex.DecodeString(hexcipher)
	if err != nil {
		return err
	}
	password, err := hex.DecodeString(hexpassword)
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
		saltsize:  16,
		saltsfile: filepath.Join(privroot, "salts"),
	}
}
