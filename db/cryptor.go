package db

// Cryptor is simple utility which manages the loading and saving of encrypted
// files. All it needs is two secure and secret keys: MAC key and an encryption
// key

// I'm not sure at all how files should be encrypted, which is partly why I'm
// building this app. Currently, we are using AES, because it's the standard,
// in CBC (cipher block chanining) mode because (1) recommends it.
// (1) https://cromwell-intl.com/cybersecurity/cipher-selection.html

// cryptor could (should) be private. I've just bounced between private and
// public so many times now, I'm don't wanna change it...

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
)

// ErrDifferentMACSum is returned if the mac sum of the file doesn't match
// the computed mac sum. This means that the file has been corrupted, be very
// careful, someone might be trying to attack
var ErrDifferentMACSum = errors.New("mac sum don't match")

// ErrInvalidPadding is returned if the padding of the plaintext during
// decryption doesn't match the expected format. This isn't a good sign, and
// may reveal to be an attack (see padding oracle)
var ErrInvalidPadding = errors.New("invalid padding")

// manages decrypting and encrypting from/to a file

// Cryptor is a simple API which writes and read encrypted files using the
// password given to the constructor
type Cryptor struct {
	block cipher.Block
	mac   hash.Hash
}

// Load opens <filename>, decrypts its content, and returns it
func (c *Cryptor) Load(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %s", err)
	}

	givenMACSum := content[:c.mac.Size()]
	ciphertext := content[c.mac.Size():]

	// check MAC sum
	c.mac.Write(ciphertext)
	computedSum := c.mac.Sum(nil)

	if hmac.Equal(computedSum, givenMACSum) {
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

	for i := len(plaintext) - npadding; i < len(plaintext); i++ {
		if plaintext[i] != npaddingbyte {
			// should more detail be given in this error message?
			return nil, ErrInvalidPadding
		}
	}

	return plaintext[:len(plaintext)-npadding], nil
}

// Save encrypts plaintext and saves it to filename
func (c *Cryptor) Save(filename string, plaintext []byte) error {
	iv, err := generateiv(c.block.BlockSize())
	if err != nil {
		return err
	}
	return c.saveWithIV(filename, plaintext, iv)
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

// NewCryptor creates a new cryptor which saves/loads encrypted files using
// the mackey and the enckey
func NewCryptor(MACKey, encryptionKey []byte) (*Cryptor, error) {

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	return &Cryptor{
		block: block,
		mac:   hmac.New(sha256.New, MACKey),
	}, nil
}

func generateiv(size int) ([]byte, error) {
	iv := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("generating IV: %s", err)
	}
	return iv, nil
}
