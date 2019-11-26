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
)

const saltsfile = "salts"

var ErrNoSaltsFile = errors.New("no saltsfile")

// FIXME: use a map internally, and then expose a static, compile-time checked
// interface

// Salts containts the *decrypted* salts
type Salts struct {
	Cipher, Password []byte
}

func GenerateNewSalts() (*Salts, error) {
	generateNewSalt := func() ([]byte, error) {
		const saltsize = 16
		salt := make([]byte, saltsize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, fmt.Errorf("generating salt: %s", err)
		}
		return salt, nil
	}

	var err error
	// Wait does that work? no nil pointer?
	var salts *Salts
	salts.Cipher, err = generateNewSalt()
	if err != nil {
		return nil, err
	}
	salts.Password, err = generateNewSalt()
	if err != nil {
		return nil, err
	}

	content := []byte(fmt.Sprintf("%s\n%s\n", salts.Cipher, salts.Password))
	if err := ioutil.WriteFile(saltsfile, content, 0644); err != nil {
		return nil, fmt.Errorf("writing salts to file system: %s", err)
	}

	return salts, nil
}

func LoadSalts() (*Salts, error) {

	f, err := os.Open(saltsfile)

	if os.IsNotExist(err) {
		return nil, ErrNoSaltsFile
	}

	if err != nil {
		return nil, fmt.Errorf("opening saltsfile: %s", err)
	}
	reader := bufio.NewReader(f)

	hexcipher, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading cipher salt: %s", err)
	}

	hexpassword, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading password salt: %s", err)
	}

	var salts *Salts
	salts.Cipher, err = hex.DecodeString(hexcipher)
	if err != nil {
		return nil, err
	}
	salts.Password, err = hex.DecodeString(hexpassword)
	if err != nil {
		return nil, err
	}

	return salts, nil
}
