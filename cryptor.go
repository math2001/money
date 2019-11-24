package main

// I'm not sure at all how files should be encrypted, which is partly why I'm
// building this app. Currently, we are using AES, because it's the standard,
// in CBC (cipher block chanining) mode because (1) recommends it.
// (1) https://cromwell-intl.com/cybersecurity/cipher-selection.html

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/pbkdf2"
)

var ErrDifferentMACSum = errors.New("mac sum don't match")
var ErrInvalidPadding = errors.New("invalid padding")

// manages decrypting and encrypting from/to a file

type Cryptor struct {
	block cipher.Block
	mac   hash.Hash
}

func (c *Cryptor) Load(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %s", err)
	}

	given_macsum := content[:c.mac.Size()]
	ciphertext := content[c.mac.Size():]

	// check MAC sum
	c.mac.Write(ciphertext)
	computed_sum := c.mac.Sum(nil)

	if hmac.Equal(computed_sum, given_macsum) {
		return nil, ErrDifferentMACSum
	}

	if len(ciphertext) < c.block.BlockSize() {
		return nil, fmt.Errorf("ciphertext too short (file corrupted)")
	}

	iv := ciphertext[:c.block.BlockSize()]
	ciphertext = ciphertext[c.block.BlockSize():]

	if len(ciphertext)%c.block.BlockSize() != 0 {
		return nil, fmt.Errorf("ciphertext length isn't a multiple of block size (file corrupted)")
	}

	// the plaintext is the buffer with the same length as ciphertext *without*
	// the iv
	plaintext := make([]byte, len(ciphertext))

	mode := cipher.NewCBCDecrypter(c.block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// remove padding from plaintext
	npaddingbyte := plaintext[len(plaintext)-1]
	npadding := int(npaddingbyte)

	if npadding > c.mac.BlockSize() {
		// should we give more details in this error message?
		return nil, ErrInvalidPadding
	}

	fmt.Println(npadding)
	for i := len(plaintext) - npadding; i < len(plaintext); i++ {
		if plaintext[i] != npaddingbyte {
			// should more detail be given in this error message?
			return nil, ErrInvalidPadding
		}
	}

	return plaintext[:len(plaintext)-npadding], nil
}

func (c *Cryptor) saveWithIV(filename string, plaintext []byte, iv []byte) error {
	blocksize := c.block.BlockSize()

	// the number of padding bytes required
	npaddingbytes := blocksize - (len(plaintext) % blocksize)

	// we add 3 threes if we have 3 padding bytes, 1 one if we have 1 padding
	// byte, etc. It seems to be the standard way of padding
	plaintext = append(plaintext, bytes.Repeat([]byte{byte(npaddingbytes)}, npaddingbytes)...)

	ciphertext := make([]byte, blocksize+len(plaintext))

	// the IV doesn't need to be secret, just unique, so we store at the
	// beginning of the ciphertext
	copy(ciphertext[:blocksize], iv)

	mode := cipher.NewCBCEncrypter(c.block, iv)
	mode.CryptBlocks(ciphertext[blocksize:], plaintext)

	// ciphertext includes the IV and the regular blocks
	c.mac.Write(ciphertext)
	c.mac.Reset()

	signature := c.mac.Sum(nil)

	if err := ioutil.WriteFile(filename, append(signature, ciphertext...), 0644); err != nil {
		return fmt.Errorf("writing file: %s", err)
	}

	return nil
}

// Encrypts plaintext and saves it to filename
func (c *Cryptor) Save(filename string, plaintext []byte) error {
	iv, err := generateiv(c.block.BlockSize())
	if err != nil {
		return err
	}
	return c.saveWithIV(filename, plaintext, iv)
}

func NewCryptor(password []byte) (*Cryptor, error) {
	// generate key from password

	salt, err := hex.DecodeString(hexsalt)
	if err != nil {
		// FIXME: can the salt be exposed in this error?
		return nil, fmt.Errorf("loading salt: %s", err)
	}
	key := pbkdf2.Key(password, salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	return &Cryptor{
		block: block,
		// FIXME: the mac key should be different than the password would it
		// even be good if it was generated *from* the password, or does it
		// have to be completely independent?
		mac: hmac.New(sha256.New, password),
	}, nil
}

func generateiv(size int) ([]byte, error) {
	iv := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("generating IV: %s", err)
	}
	return iv, nil
}
