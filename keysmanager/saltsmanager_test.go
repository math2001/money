package keysmanager_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/math2001/money/keysmanager"
)

func TestSaltsManagerBasic(t *testing.T) {
	const keylen = 32
	const (
		a = iota
		b
		c
		n
	)

	filename := getTemporaryPath(t, t.Name()+"-")
	defer os.Remove(filename)

	sm1 := keysmanager.NewSaltsManager(n, filename, keylen)

	if err := sm1.GenerateNew(); err != nil {
		t.Fatalf("generating new salts: %s", err)
	}

	sm2 := keysmanager.NewSaltsManager(n, filename, keylen)

	if err := sm2.Load(); err != nil {
		t.Fatalf("loading salts: %s", err)
	}

	if !bytes.Equal(sm1.Get(a), sm2.Get(a)) {
		t.Errorf("salt %d are different", a)
	}

	if !bytes.Equal(sm1.Get(b), sm2.Get(b)) {
		t.Errorf("salt %d are different", b)
	}

	if !bytes.Equal(sm1.Get(c), sm2.Get(c)) {
		t.Errorf("salt %d are different", c)
	}
}

func getTemporaryPath(t *testing.T, prefix string) string {
	// FIXME: loop while result path exists, with limit
	// FIXME: find a way to make that function available to everyone
	tempdir := os.TempDir()
	random := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, random); err != nil {
		t.Fatalf("generating random name: %s", err)
	}
	return filepath.Join(tempdir, prefix+hex.EncodeToString(random))
}
