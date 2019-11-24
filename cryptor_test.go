package main

import (
	"bytes"
	"os"
	"testing"
)

// FIXME: use virtual file system
// FIXME: use some kind of pre-test/post-test to setup/teardown the folder
func TestCryptor(t *testing.T) {
	password := []byte("hello world")

	cryptor, err := NewCryptor(password)
	if err != nil {
		t.Errorf("creating cryptor: %s", err)
	}

	if err := os.MkdirAll("test-store", os.ModePerm); err != nil {
		t.Fatalf("makdir test-store: %s", err)
	}

	input := []byte("aaaaa bbbb cccc dddd")

	if err := cryptor.Save("test-store/test1", input); err != nil {
		t.Fatalf("saving to test-store/test1: %s", err)
	}

	output, err := cryptor.Load("test-store/test1")
	if err != nil {
		t.Fatalf("loading from test-store/test1: %s", err)
	}

	if !bytes.Equal(input, output) {
		t.Errorf("Save(n); m = Load(); n should equal m\n%q\n%q", input, output)
	}

	if err := os.RemoveAll("test-store"); err != nil {
		t.Fatalf("remove test-store: %s", err)
	}
}
