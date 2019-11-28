package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const saltsfile = "salts"

// ErrNoSaltsFile is returned when the salt file isn't found in the private
// directory
var ErrNoSaltsFile = fmt.Errorf("no salts file (%w)", ErrPrivCorrupted)

// Salts containts the *decrypted* salts
type Salts struct {
	Cipher, Password []byte
}

func generateNewSalts(privroot string) (*Salts, error) {
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
	salts := &Salts{}
	salts.Cipher, err = generateNewSalt()
	if err != nil {
		return nil, err
	}
	salts.Password, err = generateNewSalt()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(privroot, saltsfile)
	content := []byte(fmt.Sprintf("%x\n%x\n", salts.Cipher, salts.Password))
	if err := ioutil.WriteFile(path, content, 0644); err != nil {
		return nil, fmt.Errorf("writing salts to file system: %s", err)
	}

	return salts, nil
}

func loadSalts(privroot string) (*Salts, error) {

	f, err := os.Open(filepath.Join(privroot, saltsfile))

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

	salts := &Salts{}
	salts.Cipher, err = hex.DecodeString(hexcipher[:len(hexcipher)-1])
	if err != nil {
		return nil, err
	}
	salts.Password, err = hex.DecodeString(hexpassword[:len(hexpassword)-1])
	if err != nil {
		return nil, err
	}

	return salts, nil
}
