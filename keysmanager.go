package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

// generates and loads the keys from keys.priv

var ErrKeysfileExists = errors.New("keysfile already exists")
var ErrNoKeysfile = errors.New("no keysfile")

type KeysManager struct {
	block    cipher.Block
	keysfile string
	keysize  int
}

type Keys struct {
	MAC, Encryption []byte
}

// String prevents someone printing keys without realizing that they secret. If
// he *really* wants to see the keys, he has to print them manually
// (fmt.Println(keys.MAC))
func (k Keys) String() string {
	return fmt.Sprintf("[secret!] Keys{}")
}

func (km *KeysManager) LoadKeys() (Keys, error) {
	var keys Keys

	f, err := os.Open(km.keysfile)
	if os.IsNotExist(err) {
		return keys, ErrNoKeysfile
	}
	if err != nil {
		return keys, fmt.Errorf("opening keys file: %w", err)
	}
	scanner := bufio.NewScanner(f)
	counter := 0

	// FIXME: can this key loading process be a bit smoother? Please don't be
	// too smart about it. I'm not sure there is a way
	for scanner.Scan() && counter < 2 {
		line := scanner.Text()
		if line[0] == '#' || len(line) == 0 {
			continue
		}
		keyEncrypted, err := hex.DecodeString(line)
		if err != nil {
			return keys, fmt.Errorf("Decode hex key: %s", err)
		}
		decryptedKey, err := km.decryptkey(keyEncrypted)
		if err != nil {
			return keys, fmt.Errorf("decrypting key: %s", err)
		}
		if counter == 0 {
			keys.MAC = decryptedKey
		} else if counter == 1 {
			keys.Encryption = decryptedKey
		}
		counter++
	}
	if err := scanner.Err(); err != nil {
		return keys, fmt.Errorf("scanning keys file: %s", err)
	}
	return keys, nil
}

func (km *KeysManager) decryptkey(ciphertext []byte) ([]byte, error) {
	if len(ciphertext)%km.block.BlockSize() != 0 {
		// CHECKME: can I safely give more information here?
		return nil, fmt.Errorf("length should be a multiple of blocksize")
	}

	iv := ciphertext[:km.block.BlockSize()]
	keyEncrypted := ciphertext[km.block.BlockSize():]

	keyDecrypted := make([]byte, len(keyEncrypted))

	mode := cipher.NewCBCDecrypter(km.block, iv)
	mode.CryptBlocks(keyDecrypted, keyEncrypted)

	if len(keyDecrypted) != km.keysize {
		return nil, fmt.Errorf("length should be %d, got %d", km.keysize, len(keyDecrypted))
	}

	return keyDecrypted, nil
}

func (km *KeysManager) generateHexKey() (string, error) {
	decryptedKey := make([]byte, km.keysize)

	if _, err := io.ReadFull(rand.Reader, decryptedKey); err != nil {
		return "", fmt.Errorf("generating key: %s", err)
	}

	encryptedKey := make([]byte, km.keysize+km.block.BlockSize())

	iv := encryptedKey[:km.block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("generating iv: %s", err)
	}

	mode := cipher.NewCBCEncrypter(km.block, iv)
	mode.CryptBlocks(encryptedKey[len(iv):], decryptedKey)

	return hex.EncodeToString(encryptedKey), nil
}

// GenerateKeys generates the mac and encryption keys, encrypts them and stores
// them in keysfile
func (km *KeysManager) GenerateNewKeys(password []byte) error {
	if _, err := os.Stat(km.keysfile); err == nil {
		return ErrKeysfileExists
	}

	salt, err := hex.DecodeString(hexsalt)
	if err != nil {
		// CHECKME: can the salt be exposed in this error?
		return fmt.Errorf("loading salt: %s", err)
	}
	cipherkey := pbkdf2.Key(password, salt, 4096, 32, sha256.New)

	km.block, err = aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("initiating new cipher: %s", err)
	}

	MACKey, err := km.generateHexKey()
	if err != nil {
		return err
	}
	encryptionKey, err := km.generateHexKey()
	if err != nil {
		return err
	}

	content := []byte(fmt.Sprintf("%s\n%s\n", MACKey, encryptionKey))
	if err := ioutil.WriteFile(km.keysfile, content, 0644); err != nil {
		return fmt.Errorf("writing keysfile: %s", err)
	}
	return nil
}

func (km *KeysManager) RemoveKeysfile() error {
	if err := os.Remove(km.keysfile); err != nil {
		return fmt.Errorf("removing %q keys file: %s", km.keysfile, err)
	}
	return nil
}

// Creates a key manager from the raw user password
func (km *KeysManager) Login(password []byte) error {

	salt, err := hex.DecodeString(hexsalt)
	if err != nil {
		// CHECKME: can the salt be exposed in this error?
		return fmt.Errorf("loading salt: %s", err)
	}
	cipherkey := pbkdf2.Key(password, salt, 4096, 32, sha256.New)

	cipher, err := aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("creating cipher: %s", err)
	}

	km.block = cipher
	return nil
}

func NewKeysManager() *KeysManager {
	return &KeysManager{
		keysfile: "keys.priv",
		keysize:  32,
	}
}
