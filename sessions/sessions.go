// sessions provide a way to store user information on their end (in a cookie)
// the content of the cookie is signed (but not encrypted), hence it can be
// assumed that it hasn't been altered
//
// cookie session format. Each part is base64 encoded
// alg.payload.signature
// it's heavily inspired by JWT, but I was stupid enough to do something
// *slightly* different
package sessions

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrInvalidSignature = errors.New("session cookie: payload signature didn't match")

var ErrNoSession = errors.New("no session cookie")

// S should be instantiated once, and then used for every
// request
type S struct {
	cookieName string
	// only one support signature algorithm, sha256
	alg string

	// key is the program's secret that's used to sign the cookies payload
	// IF THIS KEY IS LEAKED, ANYBODY COULD CLAIM TO BE ANYBODY ELSE EXTREMELY
	// EASILY
	key []byte
}

func NewS(config *Config) (*S, error) {
	cn := config.CookieName
	if cn == "" {
		cn = "session"
	}
	alg := config.Algorithm
	if alg == "" {
		alg = "sha256"
	}

	if alg != "sha256" {
		return nil, fmt.Errorf("unsuported algorithm: %q", alg)
	}

	if len(config.Key) == 0 {
		return nil, fmt.Errorf("empty signature key")
	}

	return &S{
		cookieName: cn,
		alg:        alg,
		key:        config.Key,
	}, nil
}

type Config struct {
	CookieName string
	Algorithm  string
	Key        []byte
}

type secret []byte

func (secret) String() string {
	return "[secret]"
}

func (s *S) Save(w http.ResponseWriter, obj interface{}) error {
	payload, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshaling obj: %s", err)
	}

	var token bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &token)
	if _, err := encoder.Write([]byte(s.alg)); err != nil {
		return fmt.Errorf("alg: %s", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("closing (flushing) partial: %s", err)
	}

	token.Write([]byte("."))
	if _, err := encoder.Write(payload); err != nil {
		return fmt.Errorf("payload: %s", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("closing (flushing) partial: %s", err)
	}
	token.Write([]byte("."))

	h := hmac.New(sha256.New, s.key)
	h.Write(token.Bytes()) // sign the whole head.base64payload. (note the extra dot)
	if _, err := encoder.Write(h.Sum(nil)); err != nil {
		return fmt.Errorf("signature: %s", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("closing (flushing) partial: %s", err)
	}

	// save bytes to cookie
	http.SetCookie(w, &http.Cookie{
		Name:  s.cookieName,
		Value: token.String(),

		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
	})

	return nil
}

// Load returns the current session. Errors: ErrInvalidSignature, ErrNoSession
func (s *S) Load(r *http.Request, dst interface{}) error {
	cookie, err := r.Cookie(s.cookieName)
	if err == http.ErrNoCookie {
		return ErrNoSession
	}
	splits := strings.Split(cookie.Value, ".")
	if len(splits) != 3 {
		return fmt.Errorf("session cookie: wrong splits")
	}

	alg, err := base64.StdEncoding.DecodeString(splits[0])
	if err != nil {
		return fmt.Errorf("session coookie: reading alg: %s", err)
	}

	if string(alg) != s.alg {
		return fmt.Errorf("session cookie: unknown alg: %q", alg)
	}

	payload, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		return fmt.Errorf("session cookie: reading payload: %s", err)
	}

	signature, err := base64.StdEncoding.DecodeString(splits[2])
	if err != nil {
		return fmt.Errorf("session cookie: reading payload: %s", err)
	}

	h := hmac.New(sha256.New, s.key)
	h.Write([]byte(fmt.Sprintf("%s.%s.", splits[0], splits[1])))

	if !hmac.Equal(h.Sum(nil), signature) {
		return ErrInvalidSignature
	}

	if err := json.Unmarshal(payload, dst); err != nil {
		return fmt.Errorf("session cookie: unmarshaling: %s", err)
	}
	return nil
}

// Remove removes the session cookie
func (s *S) Remove(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  s.cookieName,
		Value: "",

		MaxAge: 0,
	})
}
