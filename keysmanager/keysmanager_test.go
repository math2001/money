package keysmanager

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// the reason we create a different privroot folder for each test is to allow
// each of them to run concurrently

func TestKeysManagerBasic(t *testing.T) {
	t.Parallel()
	privroot := getTemporaryPath(t, "test-priv-"+t.Name())
	km := NewKeysManager(privroot)
	defer func() {
		if err := km.RemovePrivroot(); err != nil {
			t.Fatalf("tearing down (remove privroot): %s", err)
		}
	}()

	password := []byte("some random password")
	if err := km.SignUp(password); err != nil {
		t.Fatalf("signing up: %s", err)
	}

	keys1, err := km.LoadKeys()
	if err != nil {
		t.Fatalf("loading keys after sign up: %s", err)
	}

	km2 := NewKeysManager(privroot)
	if err := km2.Login(password); err != nil {
		t.Fatalf("logging in: %s", err)
	}

	keys2, err := km2.LoadKeys()
	if err != nil {
		t.Fatalf("loading keys after log in: %s", err)
	}

	if !keys1.Equal(keys2) {
		t.Fatalf("keys loaded after sign up and after log in are different")
	}
}

func TestKeysManagerWrongPassword(t *testing.T) {
	t.Parallel()
	privroot := getTemporaryPath(t, "test-priv-"+t.Name())

	km1 := NewKeysManager(privroot)
	defer func() {
		if err := km1.RemovePrivroot(); err != nil {
			t.Fatalf("tearing down (remove privroot): %s", err)
		}
	}()

	password := []byte("setup password")
	if err := km1.SignUp(password); err != nil {
		t.Fatalf("signing up: %s", err)
	}

	km2 := NewKeysManager(privroot)
	password = []byte("the wrong password!")
	err := km2.Login(password)
	if err != ErrWrongPassword {
		t.Fatalf("should have ErrorWrongPassword, got %s", err)
	}
}

func TestKeysManagerMultipleSignUp(t *testing.T) {
	t.Parallel()

	privroot := getTemporaryPath(t, "test-priv-"+t.Name())

	km1 := NewKeysManager(privroot)
	defer func() {
		if err := km1.RemovePrivroot(); err != nil {
			t.Fatalf("tearing down (remove privroot): %s", err)
		}
	}()

	password := []byte("my awesome password")
	if err := km1.SignUp(password); err != nil {
		t.Fatalf("signing up: %s", err)
	}

	err := km1.SignUp(password)
	if !errors.Is(err, ErrAlreadyLoaded) {
		t.Fatalf("should have ErrAlreadyLoggedIn, got %s", err)
	}
}

func TestKeysManagerSignUpThenLogin(t *testing.T) {
	t.Parallel()

	privroot := getTemporaryPath(t, "test-priv-"+t.Name())

	km1 := NewKeysManager(privroot)
	defer func() {
		if err := km1.RemovePrivroot(); err != nil {
			t.Fatalf("tearing down (remove privroot): %s", err)
		}
	}()

	password := []byte("my awesome password")
	if err := km1.SignUp(password); err != nil {
		t.Fatalf("signing up: %s", err)
	}

	if err := km1.Login(password); !errors.Is(err, ErrAlreadyLoaded) {
		t.Fatalf("should have ErrAlreadyLoggedIn, got %s", err)
	}
}

func TestKeysManagerMultipleLogin(t *testing.T) {
	t.Parallel()

	privroot := getTemporaryPath(t, "test-priv-"+t.Name())

	km1 := NewKeysManager(privroot)
	defer func() {
		if err := km1.RemovePrivroot(); err != nil {
			t.Fatalf("tearing down (remove privroot): %s", err)
		}
	}()

	password := []byte("my awesome password")
	if err := km1.SignUp(password); err != nil {
		t.Fatalf("signing up: %s", err)
	}

	km2 := NewKeysManager(privroot)
	if err := km2.Login(password); err != nil {
		t.Fatalf("login in: %s", err)
	}

	if err := km2.Login(password); !errors.Is(err, ErrAlreadyLoaded) {
		t.Fatalf("should have ErrAlreadyLoggedIn, got %s", err)
	}
}

// alters the password hash file, and make sure we can't log in afterwards.
// This isn't the best test because it relies on knowing how the internal
// works. However, it wouldn't make sense to expose those internals in the
// keysmanager from the perspective of someone *using* km... I guess that's
// what you get for trying to test a hack
func TestKeysManagerCorruptPasswordHashFile(t *testing.T) {
	t.Parallel()
	privroot := getTemporaryPath(t, "test-priv-"+t.Name())
	km1 := NewKeysManager(privroot)
	defer func() {
		if err := km1.RemovePrivroot(); err != nil {
			t.Fatalf("tearing down (remove privroot): %s", err)
		}
	}()

	password := []byte("some random password")
	if err := km1.SignUp(password); err != nil {
		t.Fatalf("signing up: %s", err)
	}

	_, err := km1.LoadKeys()
	if err != nil {
		t.Fatalf("loading keys after sign up: %s", err)
	}

	content, err := ioutil.ReadFile(km1.passwordhashfile)
	if err != nil {
		t.Fatalf("[hacking internals] opening passwordhashfile %q: %s", km1.passwordhashfile, err)
	}
	// this will cause the hex decoding will fail, because 'g' isn't a valid
	// hex character
	content[2] = 'g'
	if err := ioutil.WriteFile(km1.passwordhashfile, content, 0644); err != nil {
		t.Fatalf("[hacking internals] writing passwordhashfile %q: %s", km1.passwordhashfile, err)
	}

	km2 := NewKeysManager(privroot)
	err = km2.Login(password)
	if err == nil {
		t.Fatalf("we corrupted the passwordhashfile, should not be able to login")
	} else if !errors.Is(err, ErrPrivCorrupted) {
		t.Fatalf("error %q should wrap ErrPrivCorrupted", err)
	}

}

func getTemporaryPath(t *testing.T, prefix string) string {
	// FIXME: loop while result path exists, with limit
	tempdir := os.TempDir()
	random := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, random); err != nil {
		t.Fatalf("generating random name: %s", err)
	}
	return filepath.Join(tempdir, prefix+hex.EncodeToString(random))
}
