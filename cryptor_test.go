package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// FIXME: use virtual file system
func TestMain(m *testing.M) {
	if err := os.MkdirAll("test-store", os.ModePerm); err != nil {
		log.Fatalf("makdir test-store: %s", err)
	}
	code := m.Run()
	if err := os.RemoveAll("test-store"); err != nil {
		log.Fatalf("remove test-store: %s", err)
	}
	os.Exit(code)
}

func TestBasicCryptor(t *testing.T) {
	password := []byte("hello world")

	cryptor, err := NewCryptor(password)
	if err != nil {
		t.Fatalf("creating cryptor: %s", err)
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
		t.Errorf("input should equal output\n%q\n%q", input, output)
	}
}

// Test saving with one cryptor and loading with an other
func TestDifferentCryptor(t *testing.T) {
	password := []byte("some other password")

	cryptorw, err := NewCryptor(password)
	if err != nil {
		t.Fatalf("creating writting cryptor: %s", err)
	}

	input := []byte("asdf poijwqefad asdf owqiejfasldfkw")

	if err := cryptorw.Save("test-store/test2", input); err != nil {
		t.Fatalf("saving to test-store/test2: %s", err)
	}

	cryptorr, err := NewCryptor(password)
	if err != nil {
		t.Fatalf("create reading cryptor: %s", err)
	}

	output, err := cryptorr.Load("test-store/test2")
	if err != nil {
		t.Fatalf("loading from test-store/test2: %s", err)
	}

	if !bytes.Equal(input, output) {
		t.Errorf("input should equal output\n%q\n%q", input, output)
	}
}
