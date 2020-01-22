package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Entry struct {
	ID int
	Name string
	Description string
	Amount int
	Date int
	Matched bool
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/load", h(func(r *http.Request) *resp {
		f, err := os.Open("./_simpleserver/data/entries.json")
		if err != nil {
			return &resp{err: err}
		}
		decoder := json.NewDecoder(f)
		entries := make([]*Entry, 0)
		for {
			var entry *Entry
			if err := decoder.Decode(&entry); err == io.EOF {
				break
			} else if err != nil {
				return &resp{err: err}
			}
			entries = append(entries, entry)
		}
		return &resp{
			code: http.StatusOK,
			body: kv{
				"kind": "success",
				"entries": entries,
			},
		}
	})).Methods(http.MethodGet)

	api.HandleFunc("/add", h(func(r *http.Request) *resp {
		f, err := os.OpenFile("./_simpleserver/data/entries.json", os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return &resp{err: err}
		}
		var entry *Entry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			return &resp{err: err}
		}
		r.Body.Close()

		encoder := json.NewEncoder(f)
		if err := encoder.Encode(entry); err != nil {
			return &resp{err: err}
		}

		return &resp{
			code: http.StatusOK,
			body: kv{
				"kind": "success",
			},
		}
	})).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe("localhost:9999", router))
}

type kv map[string]interface{}

type resp struct {
	code   int
	body   kv
	reader io.Reader
	err    error
}

func h(handler func(*http.Request) *resp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := handler(r)
		if response == nil {
			panic("nil response")
		}
		if response.reader != nil && response.body != nil {
			response.err = errors.New("got reader and body")
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		if response.err != nil {
			log.Printf("internal error: %s", response.err)
			internalError(w)
			return
		}

		if (response.code == 0) {
			log.Printf("response code is 0")
			internalError(w)
			return
		}
		w.WriteHeader(response.code)

		if response.err != nil {
			panic("response err != nil")
		}

		if response.reader != nil {
			if _, err := io.Copy(w, response.reader); err != nil {
				panic(err)
			}
		}

		if response.body != nil {
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(response.body); err != nil {
				panic(err)
			}
		}
	}
}

func internalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := fmt.Fprintf(w, `{"kind": "internal error"}`); err != nil {
		panic(err)
	}
}