package main

import (
	"bufio"
	"bytes"
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
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

// ErrKeysfileExists is returned if the keysfile already exists
var ErrKeysfileExists = errors.New("keysfile already exists")

// ErrNoKeysfile is returned if there is no keysfile to read the keys from
var ErrNoKeysfile = errors.New("no keysfile")

// ErrWrongPassword is returned when the hash store hashpasswordfile doesn't
// match with the hash of the typed password. See Login function
var ErrWrongPassword = errors.New("wrong password")

// ErrAlreadyLoggedIn is returned by Login when the user tries to log in
// multiple twice
var ErrAlreadyLoggedIn = errors.New("already logged in")

// KeysManager loads the different keys from a file (keysfile) and decrypts
// them using the password. It can also generate new keys in place of the old
// ones
type KeysManager struct {
	block                       cipher.Block
	privroot, keysfile, saltdir string
	passwordhashfile            string
	keysize                     int
	salts                       *Salts
}

// Keys contains the *decrypted* keys. Plaintext. Be careful
type Keys struct {
	MAC, Encryption []byte
}

// String prevents someone printing keys without realizing that they secret. If
// he *really* wants to see the keys, he has to print them manually
// (fmt.Println(keys.MAC))
func (k Keys) String() string {
	return fmt.Sprintf("[secret!] Keys{}")
}

// LoadKeys loads the keys from the file
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
		decryptedKey, err := km.decryptKey(keyEncrypted)
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

func (km *KeysManager) decryptKey(ciphertext []byte) ([]byte, error) {
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

// encryptKey encrypts the given key with the current cipher, and then hex
// encodes it
func (km *KeysManager) encryptKey(decryptedKey []byte) (string, error) {
	encryptedKey := make([]byte, len(decryptedKey)+km.block.BlockSize())

	iv := encryptedKey[:km.block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("generating iv to encrypt key: %s", err)
	}

	mode := cipher.NewCBCEncrypter(km.block, iv)
	mode.CryptBlocks(encryptedKey[len(iv):], decryptedKey)

	return hex.EncodeToString(encryptedKey), nil
}

// generateHexKey create a new key, encrypts it with the current block cipher,
// and then hex encodes it
// FIXME: this is a really bad name, we don't understand that it is being
// securely encrypted from it.
func (km *KeysManager) generateHexKey() (string, error) {
	decryptedKey := make([]byte, km.keysize)

	if _, err := io.ReadFull(rand.Reader, decryptedKey); err != nil {
		return "", fmt.Errorf("generating key: %s", err)
	}

	return km.encryptKey(decryptedKey)
}

// GenerateNewKeys generates the mac and encryption keys, encrypts them and
// stores them in keysfile
// FIXME: this function shouldn't be exposed...
func (km *KeysManager) GenerateNewKeys(password []byte) error {
	if _, err := os.Stat(km.keysfile); err == nil {
		return ErrKeysfileExists
	}

	cipherkey := pbkdf2.Key(password, km.salts.Cipher, 4096, 32, sha256.New)

	cipher, err := aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("initiating new cipher: %s", err)
	}
	// do it in two steps to not erase km.block in case there is an error
	km.block = cipher

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

// RemoveKeysfile deletes the current keysfile. The keys stored will be
// permanantely lost
func (km *KeysManager) RemoveKeysfile() error {
	if err := os.Remove(km.keysfile); err != nil {
		return fmt.Errorf("removing %q keys file: %s", km.keysfile, err)
	}
	return nil
}

// ChangePassword changes the user's password. Note that the user must already
// be logged in. If it isn't case, use set password
func (km *KeysManager) ChangePassword(newpassword []byte) error {
	keys, err := km.LoadKeys()
	if err != nil {
		return fmt.Errorf("loading keys: %s", err)
	}

	// replace cipher(currentpassword) with a cipher(newpassword)

	cipherkey := pbkdf2.Key(newpassword, km.salts.Cipher, 4096, 32, sha256.New)

	cipher, err := aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("creating cipher: %s", err)
	}
	// we do this in two steps (cipher = Cipher() then km.block = cipher) so
	// that we don't set km.block to nil if aes.NewCipher fails. It makes sense
	// because if we fail here, the only action that fails is the ChangePassword
	// action, not the entire app.
	km.block = cipher

	// FIXME: how can we make sure here that all the keys are updated?
	// (ie. for example if we add a field to Keys, how can we write some code
	// that will generate an error if we have a third field in Keys but only update two here)
	MACKey, err := km.encryptKey(keys.MAC)
	if err != nil {
		return fmt.Errorf("encrypting keys with new password: %s", err)
	}
	encryptionKey, err := km.encryptKey(keys.MAC)
	if err != nil {
		return fmt.Errorf("encrypting keys with new password: %s", err)
	}

	// encrypt the current keys with the new password
	content := []byte(fmt.Sprintf("%s\n%s\n", MACKey, encryptionKey))
	if err := ioutil.WriteFile(km.keysfile, content, 0644); err != nil {
		return fmt.Errorf("writing keysfile: %s", err)
	}
	return nil
}

// SignUp makes the priv directory, creates the salts, password hash file and
// keys
func (km *KeysManager) SignUp(password []byte) error {

	if err := os.MkdirAll(km.privroot, 0755); err != nil {
		return fmt.Errorf("making privroot directory %q: %s", km.privroot, err)
	}

	var err error
	km.salts, err = GenerateNewSalts()
	if err != nil {
		return fmt.Errorf("generating salts: %s", err)
	}

	passwordhash := pbkdf2.Key(password, km.salts.Password, 4096, 32, sha256.New)
	if err := ioutil.WriteFile(km.passwordhashfile, []byte(hex.EncodeToString(passwordhash)), 0644); err != nil {
		return fmt.Errorf("writing password hash to file: %s", err)
	}

	if err := km.GenerateNewKeys(password); err != nil {
		return fmt.Errorf("generating new keys: %s", err)
	}

	cipherkey := pbkdf2.Key(password, km.salts.Cipher, 4096, 32, sha256.New)

	cipher, err := aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("creating cipher: %s", err)
	}

	km.block = cipher

	return nil
}

// FIXME: how do you transfer keys? Just copy paste the priv/ folder?

// Login creates the block cipher from the password, which will then be used
// to decrypt the keys from the file
func (km *KeysManager) Login(password []byte) error {

	if km.block != nil {
		return ErrAlreadyLoggedIn
	}

	var err error
	km.salts, err = LoadSalts()
	if err != nil {
		return fmt.Errorf("loading salts: %w", err)
	}

	// check that the password's hash matches the hash stored in passwordhashfile
	// this is just extra feature to error as early as possible if the password
	// is wrong. We could not do this step, and let the user log in with a
	// wrong password: he decrypt his keys wrongly, and hence wouldn't be able
	// to retrieve the original content from his files (padding error, or if
	// he is lucky, gibberish)

	hexpasswordhash, err := ioutil.ReadFile(km.passwordhashfile)
	if os.IsNotExist(err) {
		fmt.Println("No hash file to compare password against.")
		// FIXME: what should we do here? Store the new password in hashfile
		// and keep going? What if it's wrong? next time the user tries to log
		// in with the right password, they'll be kicked out! Maybe we should
		// have a second command, like login-force, which wouldn't have this
		// checking feature
		return fmt.Errorf("not implemented")
	}
	if err != nil {
		return fmt.Errorf("reading from password hash file: %s", err)
	}

	passwordhash, err := hex.DecodeString(string(hexpasswordhash))
	if err != nil {
		return fmt.Errorf("decoding hex password hash: %s", err)
	}

	if !bytes.Equal(passwordhash, pbkdf2.Key(password, km.salts.Password, 4096, 32, sha256.New)) {
		// that doesn't *actually* mean that the password is wrong. It could
		// be that the stored hash is wrong, the that this password would sill
		// succesfully decode the files. This however shouldn't happen, check
		// FIXME above
		return ErrWrongPassword
	}

	cipherkey := pbkdf2.Key(password, km.salts.Cipher, 4096, 32, sha256.New)

	cipher, err := aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("creating cipher: %s", err)
	}

	km.block = cipher
	return nil
}

// NewKeysManager create a new KeysManager with some sane default
func NewKeysManager() (*KeysManager, error) {

	// the files in the priv folder don't *have* to be secret. Technically, you
	// could share it with everyone, and this application should still be
	// secure. However, it would make an attacker's job easier if he had
	// acccess to them.

	// keys.priv: the keys are encrypted with the user's password. Those keys
	// are what Cryptor uses to encrypt user files

	// salts: salts just have to be unique to fight against rainbow tables. If
	// an attacker had access to those files before actually breaking into the
	// application, he could generate a table of hashes using that salt to very
	// quickly find the user's password (provided the it is somewhat common)

	// passwordhash.priv: this is just a hash of the user's password. It is
	// technically secure, but best kept secret.

	km := &KeysManager{
		privroot: "priv",
		keysize:  32,
	}

	km.keysfile = filepath.Join(km.privroot, "keys.priv")
	km.saltdir = filepath.Join(km.privroot, "salts")
	km.passwordhashfile = filepath.Join(km.privroot, "passwordhash.priv")

	return km, nil
}
