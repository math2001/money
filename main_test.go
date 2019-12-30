package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/math2001/money/server"
)

// makes sure that the legal request do work
func TestLegalMethods(t *testing.T) {

	dataroot, err := ioutil.TempDir("", "money-test-"+t.Name())
	if err != nil {
		t.Fatalf("creating temporary directory: %s", err)
	}
	// remove the dataroot dir so that when the api is created, it detects that
	// it doesn't exists and initiates itself
	if err := os.Remove(dataroot); err != nil {
		t.Fatalf("removing temporary directory: %s", err)
	}

	defer func() {
		if err := os.RemoveAll("test-tmp-" + t.Name()); err != nil {
			t.Fatalf("remove temporary test dir: %s", err)
		}
	}()

	var logs strings.Builder
	logs.WriteString("Logs:\n")
	log.SetOutput(&logs)
	handler, err := server.New(dataroot, "0.0.0.0", []byte(""))

	type headers map[string]string
	type resp struct {
		code    int
		headers headers
	}
	cases := map[*http.Request]resp{
		httptest.NewRequest("GET", "/", nil): resp{
			code:    http.StatusOK,
			headers: headers{"Content-Type": "text/html; charset=utf-8"},
		},
		httptest.NewRequest("GET", "/css/main.css", nil): resp{
			code:    http.StatusOK,
			headers: headers{"Content-Type": "text/css; charset=utf-8"},
		},
		httptest.NewRequest("GET", "/js/sw.js", nil): resp{
			code:    http.StatusOK,
			headers: headers{"Content-Type": "application/javascript; charset=utf-8", "Service-Worker-Allowed": "/"},
		},
		httptest.NewRequest("GET", "/manifest.json", nil): resp{
			code:    http.StatusOK,
			headers: headers{"Content-Type": "application/json; charset=utf-8"},
		},
		httptest.NewRequest("GET", "/random", nil): resp{
			code:    http.StatusOK,
			headers: headers{"Content-Type": "text/html; charset=utf-8"},
		},
		httptest.NewRequest("GET", "/randomtrailing/", nil): resp{
			code:    http.StatusPermanentRedirect,
			headers: headers{"Content-Type": "text/html; charset=utf-8", "Location": "/randomtrailing"},
		},
		httptest.NewRequest("GET", "/random/nested", nil): resp{
			code:    http.StatusOK,
			headers: headers{"Content-Type": "text/html; charset=utf-8"},
		},
		httptest.NewRequest("GET", "/api", nil): resp{
			code:    http.StatusMovedPermanently,
			headers: headers{"Location": "/api/"},
		},
		httptest.NewRequest("GET", "/api/", nil): resp{
			code:    http.StatusOK,
			headers: headers{"Content-Type": "application/json; charset=utf-8"},
		},
		httptest.NewRequest("POST", "/api/login", nil): resp{
			headers: headers{"Content-Type": "application/json; charset=utf-8"},
		},
		httptest.NewRequest("POST", "/api/signup", nil): resp{
			headers: headers{"Content-Type": "application/json; charset=utf-8"},
		},
	}

	for req, expected := range cases {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		resp := recorder.Result()
		if expected.code != 0 && resp.StatusCode != expected.code {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("reading response body: %s", err)
			} else {
				t.Logf("%q", body)
			}
			t.Errorf("%q: status code: got %d expected %d", req.URL.Path, resp.StatusCode, expected.code)
		}
		for name, expectedValue := range expected.headers {
			if resp.Header.Get(name) != expectedValue {
				t.Errorf("%q header %q: actual %q, expected %q", req.URL.Path, name, resp.Header.Get(name), expectedValue)
			}
		}
	}
	t.Log(logs.String())
}
