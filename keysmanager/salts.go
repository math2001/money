package keysmanager

// saltmanager stores salts *in clear* in a text file. They can't be encrypted,
// have a look at https://github.com/math2001/notes/blob/4ffe5526da0a4fa0870ec3ddc80d46051327e46e/encryption/requirements-to-encrypt-text.md#encrypt-fixed-size-text
// it's within the keysmanager, because no one would need to store random hex
// numbers in a file...

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ErrNoSaltsFile is returned when the salt file isn't found in the private
// directory
var ErrNoSaltsFile = fmt.Errorf("no salts file (%w)", ErrPrivCorrupted)

type salts struct {
	cipher, password []byte
}

type saltsManager struct {
	// path to the file which will store the salts
	saltsfile string
	// the saltsize, in byte
	saltsize int
	// the salts
	salts *salts
}

func newSaltsManager(saltsfile string, saltsize int) *saltsManager {
	return &saltsManager{
		saltsfile: saltsfile,
		saltsize:  saltsize,
		salts:     &salts{},
	}
}

func (sm *saltsManager) generateNewSalts() error {
	generateNewSalt := func() ([]byte, error) {
		salt := make([]byte, sm.saltsize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, fmt.Errorf("generating salt: %s", err)
		}
		return salt, nil
	}

	var err error
	sm.salts.cipher, err = generateNewSalt()
	if err != nil {
		return err
	}
	sm.salts.password, err = generateNewSalt()
	if err != nil {
		return err
	}

	content := []byte(fmt.Sprintf("%x\n%x\n", sm.salts.cipher, sm.salts.password))
	if err := ioutil.WriteFile(sm.saltsfile, content, 0644); err != nil {
		return fmt.Errorf("writing salts to file system: %s", err)
	}

	return nil
}

func (sm *saltsManager) loadSalts() error {
	if sm.salts == (&salts{}) {
		return fmt.Errorf("already loaded salts (%w)", ErrAlreadyLoaded)
	}

	f, err := os.Open(sm.saltsfile)
	if err != nil {
		return fmt.Errorf("opening saltsfile: %s (%w)", err, ErrPrivCorrupted)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	hexcipher, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading cipher salt: %s", err)
	}

	hexpassword, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading password salt: %s", err)
	}

	sm.salts.cipher, err = hex.DecodeString(hexcipher[:len(hexcipher)-1])
	if err != nil {
		return err
	}

	sm.salts.password, err = hex.DecodeString(hexpassword[:len(hexpassword)-1])
	if err != nil {
		return err
	}

	return nil
}
