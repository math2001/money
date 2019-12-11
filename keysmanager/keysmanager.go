/*
keysmanager stores secret keys encrypted with a password. It creates its own
directory with the 3 different files that it needs

    keys:

        the keys encrypted with scrypt(userPassword, saltK). userPassword is
        the same for every different key, it's just the salt that changes

    salts:

        the salts, stored in clear text

    passwordhash:

        contains the password hash. It's just a utility to be able to tell
        whether the user gave the right password when he tries to log in. It
        doesn't make the system any more secure (we could let the user go
        through, and he would just wrongly decrypt the keys)
*/
package keysmanager

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

	// FIXME: please use scrypt!!
	"golang.org/x/crypto/pbkdf2"
)

// all those errors shouldn't be exported because they aren't unsed anywhere
// else in the app. *Only export when you need it please*

// ErrPrivCorrupted is tag error which indicates that the priv directory isn't
// right (missing file, altered keys, etc...)
var ErrPrivCorrupted = errors.New("priv directory corrupted")

// ErrWrongPassword is returned when the hash store hashpasswordfile doesn't
// match with the hash of the typed password. See Login function
var ErrWrongPassword = errors.New("wrong password")

// ErrAlreadyLoaded is a multi purpose tagging error used to indicate when an
// action that should have been done only once was executed mulitple times. For
// example, you will get this error if you try to login or load keys more than
// once for example
var ErrAlreadyLoaded = errors.New("already loaded")

const (
	saltCipher = iota
	saltMac
	saltPassword
)

// KeysManager loads the different keys from a file (keysfile) and decrypts
// them using the password. It can also generate new keys in place of the old
// ones
type KeysManager struct {
	block   cipher.Block
	keysize int
	sm      *SM

	privroot         string
	keysfile         string
	passwordhashfile string
}

// NewKeysManager create a new KeysManager with some sane default
func NewKeysManager(privroot string) *KeysManager {

	// the files in the priv folder don't *have* to be secret. Technically, you
	// could share it with everyone, and this application should still be
	// secure. However, it would make an attacker's job easier if he had
	// acccess to them.

	// keys: the keys are encrypted with the user's password (and the salt
	// cipherSalt). Those keys are what Cryptor uses to encrypt user files

	// salts: salts just have to be unique to fight against rainbow tables. If
	// an attacker had access to those files before actually breaking into the
	// application, he could generate a table of hashes using that salt to very
	// quickly find the user's password (provided the it is somewhat common)

	// passwordhash: this is just a hash of the user's password (uses
	// saltPassword). It is technically secure, but best kept secret.

	km := &KeysManager{
		privroot: privroot,
		keysize:  32,
	}

	km.keysfile = filepath.Join(km.privroot, "keys")
	km.passwordhashfile = filepath.Join(km.privroot, "passwordhash")

	// CHECKME: the original saltsize was 16 I think... Is that fine?
	km.sm = NewSaltsManager(2, filepath.Join(km.privroot, "salts"), 32)
	return km
}

// SignUp makes the priv directory, creates the salts, password hash file and
// generates new keys
func (km *KeysManager) SignUp(password []byte) error {

	err := os.Mkdir(km.privroot, 0755)
	if os.IsExist(err) {
		return fmt.Errorf("already signed up (%w)", ErrAlreadyLoaded)
	}
	if err != nil {
		return fmt.Errorf("making privroot directory %q: %s", km.privroot, err)
	}

	if err := km.sm.GenerateNew(); err != nil {
		return fmt.Errorf("generating salts: %s", err)
	}

	passwordhash := pbkdf2.Key(password, km.sm.Get(saltPassword), 4096, 32, sha256.New)
	if err := ioutil.WriteFile(km.passwordhashfile, []byte(hex.EncodeToString(passwordhash)), 0644); err != nil {
		return fmt.Errorf("writing password hash to file: %s", err)
	}

	if err := km.generateNewKeys(password); err != nil {
		return fmt.Errorf("generating new keys: %s", err)
	}

	cipherkey := pbkdf2.Key(password, km.sm.Get(saltCipher), 4096, 32, sha256.New)

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

	// FIXME: export this to a function IsLoggedIn? How else could we check?
	if km.block != nil {
		return fmt.Errorf("logging in (%w)", ErrAlreadyLoaded)
	}

	var err error
	if err := km.sm.Load(); err != nil {
		// wrap because this could contain the tag ErrPrivCorrupted
		return fmt.Errorf("loading salts: %w", err)
	}

	// check that the password's hash matches the hash stored in passwordhashfile
	// this is just extra feature to error as early as possible if the password
	// is wrong. We could not do this step, and let the user log in with a
	// wrong password: he decrypt his keys wrongly, and hence wouldn't be able
	// to retrieve the original content from his files (padding error, or if
	// he is lucky, gibberish)

	hexpasswordhash, err := ioutil.ReadFile(km.passwordhashfile)
	if err != nil {
		return fmt.Errorf("reading from password hash file: %s (%w)", err, ErrPrivCorrupted)
	}

	passwordhash, err := hex.DecodeString(string(hexpasswordhash))
	if err != nil {
		return fmt.Errorf("decoding hex password hash: %s (%w)", err, ErrPrivCorrupted)
	}

	if !bytes.Equal(passwordhash, pbkdf2.Key(password, km.sm.Get(saltPassword), 4096, 32, sha256.New)) {
		// that doesn't *actually* mean that the password is wrong. It could
		// be that the stored hash is wrong, the that this password would sill
		// succesfully decode the files. This however shouldn't happen, check
		// FIXME above
		return ErrWrongPassword
	}

	cipherkey := pbkdf2.Key(password, km.sm.Get(saltCipher), 4096, 32, sha256.New)

	cipher, err := aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("creating cipher: %s", err)
	}

	km.block = cipher
	return nil
}

// HasSignedUp returns true of the privroot directory exists, even if we can't
// read from it. That's because if the user doesn't have the permission for
// example, he won't be able to create his priv directory by signing up.
// So, we let .Login report the error, because it will know best what to do
// based (whereas this function is just general-purposed)
func (km *KeysManager) HasSignedUp() bool {
	_, err := os.Stat(km.privroot)
	return err == nil || os.IsExist(err)
}

// RemovePrivroot removes permanantely the private folder. If you run that,
// you loose your keys (ie. you won't be able to decrypt your files anymore)
func (km *KeysManager) RemovePrivroot() error {
	return os.RemoveAll(km.privroot)
}

// LoadKeys loads the keys from the keys file
//
// FIXME: do something so that we ensure that we only load the keys once
func (km *KeysManager) LoadKeys() (Keys, error) {
	f, err := os.Open(filepath.Join(km.privroot, "keys"))
	if err != nil {
		return Keys{}, fmt.Errorf("opening keys file: %s (%w)", err, ErrPrivCorrupted)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	decryptKey := func(reader *bufio.Reader) ([]byte, error) {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read line from keysfile: %s", err)
		}

		// remove the line ending \n
		ciphertext, err := hex.DecodeString(line[:len(line)-1])
		if err != nil {
			return nil, fmt.Errorf("decode hex key: %s", err)
		}

		if len(ciphertext)%km.block.BlockSize() != 0 {
			return nil, fmt.Errorf("length should be a multiple of blocksize, got %d", len(ciphertext))
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

	var keys Keys
	// first it's encryption, and then mac (alphabetical order)
	keys.Encryption, err = decryptKey(reader)
	if err != nil {
		return keys, err
	}

	keys.MAC, err = decryptKey(reader)
	if err != nil {
		return keys, err
	}
	return keys, nil
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

func (km *KeysManager) generateNewKeys(password []byte) error {
	if _, err := os.Stat(km.keysfile); err == nil {
		// we don't export that error, although it's static, because nowhere
		// in the application should that happen
		return fmt.Errorf("keysfile already exists")
	}

	cipherkey := pbkdf2.Key(password, km.sm.Get(saltCipher), 4096, 32, sha256.New)

	cipher, err := aes.NewCipher(cipherkey)
	if err != nil {
		return fmt.Errorf("initiating new cipher: %s", err)
	}
	// do it in two steps to not erase km.block in case there is an error
	km.block = cipher

	generateHexKey := func() (string, error) {
		decryptedKey := make([]byte, km.keysize)

		if _, err := io.ReadFull(rand.Reader, decryptedKey); err != nil {
			return "", fmt.Errorf("generating key: %s", err)
		}

		return km.encryptKey(decryptedKey)
	}

	MACKey, err := generateHexKey()
	if err != nil {
		return err
	}
	encryptionKey, err := generateHexKey()
	if err != nil {
		return err
	}

	content := []byte(fmt.Sprintf("%s\n%s\n", MACKey, encryptionKey))
	if err := ioutil.WriteFile(km.keysfile, content, 0644); err != nil {
		return fmt.Errorf("writing keysfile: %s", err)
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

	cipherkey := pbkdf2.Key(newpassword, km.sm.Get(saltCipher), 4096, 32, sha256.New)

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
