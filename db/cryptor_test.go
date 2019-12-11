package db

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// FIXME: use virtual file system
const storedir = "test-store"

// fill with a bunch of pre-generated keys so that we can have deterministic
// keys
var keys [][]byte

func TestBasicCryptor(t *testing.T) {
	t.Parallel()
	storefile := filepath.Join(storedir, "test-"+t.Name())
	cryptor, err := NewCryptor(keys[0], keys[1])
	if err != nil {
		t.Fatalf("creating cryptor: %s", err)
	}

	input := []byte("aaaaa bbbb cccc dddd")

	if err := cryptor.Save(storefile, input); err != nil {
		t.Fatalf("saving to %s: %s", storefile, err)
	}

	output, err := cryptor.Load(storefile)
	if err != nil {
		t.Fatalf("loading from %s: %s", storefile, err)
	}

	if !bytes.Equal(input, output) {
		t.Errorf("input should equal output\n%q\n%q", input, output)
	}
}

// Test saving with one cryptor and loading with an other
func TestDifferentCryptor(t *testing.T) {
	t.Parallel()
	storefile := filepath.Join(storedir, "test-"+t.Name())
	cryptorw, err := NewCryptor(keys[2], keys[3])
	if err != nil {
		t.Fatalf("creating writting cryptor: %s", err)
	}

	input := []byte("asdf poijwqefad asdf owqiejfasldfkw")

	if err := cryptorw.Save(storefile, input); err != nil {
		t.Fatalf("saving to %s: %s", storefile, err)
	}

	cryptorr, err := NewCryptor(keys[2], keys[3])
	if err != nil {
		t.Fatalf("create reading cryptor: %s", err)
	}

	output, err := cryptorr.Load(storefile)
	if err != nil {
		t.Fatalf("loading from %s: %s", storefile, err)
	}

	if !bytes.Equal(input, output) {
		t.Errorf("input should equal output\n%q\n%q", input, output)
	}
}
func returningTestMain(m *testing.M) int {
	if err := os.MkdirAll(storedir, os.ModePerm); err != nil {
		log.Fatalf("makdir %s: %s", storedir, err)
	}
	defer func() {
		if err := os.RemoveAll(storedir); err != nil {
			log.Fatalf("remove %s: %s", storedir, err)
		}
	}()

	// $ make money-tools
	// $ money-tools generate-new-keys 20
	var hexkeys = [...]string{
		"9fca99b7144c5c2ca9e991b7cc080f2ade0b5127e34ce879004e83489907e242",
		"998a17caf56a1b1a38c3e94cb41d9c5525d5688a87c8ef62b1e9d302da94120e",
		"5614483f1097ac22b9572330cfb4c2b39840b58a746fdc927ec9002f9a861b04",
		"4f1807073e378bd71c391c6b493f5f676dc1b3a50aec0c65f1d5cbcc658c610d",
		"0ce869f8addda6a2ceaff1aa32daad7464c20b60ad955727886f879346a86e1e",
		"f47a530f6a89ed1708dccce7cdc90c0a278bf1e770d3804f0227a732cfcbcf2e",
		"70686b747b59aa268e1b46b2f57f1c8e2dbeb71fd882b35213d729b066dc7e81",
		"e1e3f825e8de7e3efe04c57953cfea75640e71e29d770e7bb23f689d5670de24",
		"890c184fb207b762c6819ec9b48a6921feaf599d1e15c37770ec41e01fcc9345",
		"43f5cfec65368d6a9aac044820b23bdfb6ed13359525c77ab0f4146b91d3cd6f",
		"dfdef00225d2d4e3c3c6b5dc64bc5cf6afea6159329e2278c8452739714515b2",
		"3a58693096429fed4cd120ebd0016fdd928c1febacb2cfb6e04c27366a5b7be0",
		"39a359d5c578c13035a6669e1f5b8f3700fce8c3c488328739d3265fe76c60dc",
		"06eaf4eaae5aa55fd8c5cf6f8ffa8f27f508f80ebbe5e6240eb8d79bdfdbdd1b",
		"f23a09f6c1acd8f5b047fa50385d67f19d16d3d7c05b5146ee782bd7e5cc3055",
		"2c5686ce477435ae2e24a57e075e2f3fae4247e1e76ed33a109f2d5b55de1a50",
		"d687009cef96b333f5f814b041f3c773b543beac9faf944d1880a9ebf44bdaf1",
		"1576b8ee14ca73827b5a3c06e41b22d9afc14473a3b00fdabc7c1b3a530632a3",
		"a5605f2921b98eb1a0487376f5459a2f1745adfd8fbe4c76a7b46f9f2aa7b83e",
		"e8728fd1f1fe3e303925f57e623dd19a56e7d8a996a709747393b0f2933fb4bb",
	}
	keys = make([][]byte, len(hexkeys))

	for i, hexkey := range hexkeys {
		key, err := hex.DecodeString(hexkey)
		if err != nil {
			fmt.Printf("[internal error] decoding pre-generated keys: %s\n", err)
			return 1 // don't log.Fatal to allow for tear down
		}
		keys[i] = key
	}
	return m.Run()
}

// use custom test main to allow returningTestMain to run its deferred functions
func TestMain(m *testing.M) {
	os.Exit(returningTestMain(m))
}
