package keysmanager

// saltmanager stores salts *in clear* in a text file. They can't be encrypted,
// have a look at https://github.com/math2001/notes/blob/4ffe5526da0a4fa0870ec3ddc80d46051327e46e/encryption/requirements-to-encrypt-text.md#encrypt-fixed-size-text
// it's within the keysmanager, because no one would need to store random hex
// numbers in a file...

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// ErrNoSaltsFile is returned when the salt file isn't found in the private
// directory
var ErrNoSaltsFile = fmt.Errorf("no salts file (%w)", ErrPrivCorrupted)

type salt []byte

func (salt) String() string {
	return "[secret!] salt"
}

type SM struct {
	// n is the number of salt the user wants to store
	n int
	// path to the file which will store the salts
	file string
	// the saltsize, in byte
	size int
	// the salts
	salts []salt
}

// NewSaltsManager will store in clear n salts of length size byte in file
func NewSaltsManager(n int, file string, size int) *SM {
	return &SM{
		n:    n,
		file: file,
		size: size,
	}
}

func (sm *SM) GenerateNew() error {
	f, err := os.Create(sm.file)
	if err != nil {
		return fmt.Errorf("creating salt file %q: %s", sm.file, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	buf.Grow(sm.size)
	w := io.MultiWriter(hex.NewEncoder(f), &buf)
	for i := 0; i < sm.n; i++ {
		if _, err := io.CopyN(w, rand.Reader, int64(sm.size)); err != nil {
			return fmt.Errorf("writing salt #%d to file %q: %s", i, sm.file, err)
		}
		f.Write([]byte("\n"))
		sm.salts = append(sm.salts, buf.Bytes())
		buf.Reset()
	}
	return nil
}

func (sm *SM) Load() error {
	if len(sm.salts) > 0 {
		return fmt.Errorf("already loaded salts (%w)", ErrAlreadyLoaded)
	}

	f, err := os.Open(sm.file)
	if err != nil {
		return fmt.Errorf("opening saltsfile: %s (%w)", err, ErrPrivCorrupted)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	sm.salts = make([]salt, sm.n)

	for i := 0; i < sm.n; i++ {
		hexsalt, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading #%d salt: %s", i, err)
		}
		plainsalt, err := hex.DecodeString(hexsalt[:len(hexsalt)-1])
		if err != nil {
			return err
		}
		sm.salts[i] = plainsalt

		// FIXME: check if we reach EOF (otherwise we are corrupted)
	}

	return nil
}

// GetSalt returns the ith salt (0 based). Panics if i >= n
func (sm *SM) Get(i int) []byte {
	// display friendlier panic
	if len(sm.salts) == 0 {
		panic("salts haven't been loaded")
	} else if i >= sm.n {
		panic(fmt.Sprintf("trying to load salt[%d], only got %d salts", i, sm.n))
	}
	return sm.salts[i]
}
