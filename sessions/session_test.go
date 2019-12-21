package sessions_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/math2001/money/sessions"
)

const cookieName = "session"

func TestNoSession(t *testing.T) {
	s := NewS(t)
	req := httptest.NewRequest(http.MethodGet, "/api/sample", nil)
	var obj interface{}
	err := s.Load(req, &obj)
	if err != nil {
		t.Fatalf("loading from empty request: %s", err)
	}
	if obj != nil {
		t.Fatalf("loading from empty request, got non-nil obj: %s", obj)
	}
}

func TestEmpty(t *testing.T) {
	s := NewS(t)

	req := httptest.NewRequest(http.MethodGet, "/api/sample", nil)
	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: "",
	})

	var obj interface{}
	err := s.Load(req, &obj)
	if err == nil {
		t.Fatalf("loading empty cookie, should have error")
	}
	if obj != nil {
		t.Fatalf("loading from empty cookie, got non-nil obj: %s", obj)
	}
}

func TestNormal(t *testing.T) {
	s := NewS(t)
	expected := map[string]string{
		"hello": "world",
	}

	w := httptest.NewRecorder()
	s.Save(w, expected)

	req := httptest.NewRequest(http.MethodGet, "/api/sample", nil)
	req.AddCookie(w.Result().Cookies()[0])

	var actual map[string]string
	err := s.Load(req, &actual)
	if err != nil {
		t.Fatalf("loading session: %s", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("original != loaded: \n%v\n%v", actual, expected)
	}
}

func TestHack(t *testing.T) {
	s := NewS(t)

	w := httptest.NewRecorder()
	s.Save(w, map[string]interface{}{
		"id": 2,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/sample", nil)
	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: craftSession(t),
	})

	err := s.Load(req, nil)
	if err == nil {
		t.Fatalf("hacked session! should have error")
	}

	// horrible error checking, I will test the warning system that will
	// be implemented later on
	if !strings.Contains(err.Error(), "signature") {
		t.Fatalf("err should have been about signature not matching, got %s", err)
	}

}

func NewS(t *testing.T) *sessions.S {
	t.Helper()
	key := make([]byte, 32)
	seed := time.Now().UnixNano()
	t.Logf("seed: %d\n", seed)
	if _, err := io.ReadFull(rand.New(rand.NewSource(seed)), key); err != nil {
		t.Fatalf("generating deterministic key: %s", err)
	}
	s, err := sessions.NewS(&sessions.Config{
		Key:        key,
		CookieName: cookieName,
	})
	if err != nil {
		t.Fatalf("creating S: %s", err)
	}
	return s
}

func craftSession(t *testing.T) string {
	t.Helper()

	a := base64.StdEncoding.EncodeToString([]byte("sha256"))
	b := base64.StdEncoding.EncodeToString([]byte(`{"id":3}`))
	h := hmac.New(sha256.New, []byte("random key..."))
	h.Write([]byte(fmt.Sprintf("%s.%s", a, b)))
	c := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s.%s.%s", a, b, c)
}
