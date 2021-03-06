package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/math2001/money/api"
	"github.com/math2001/money/db"
	"github.com/math2001/money/keysmanager"
	"github.com/math2001/money/sessions"
)

// Session is the content of the session cookie
type Session struct {
	ID       int
	Email    string
	Password secret
}

type secret []byte

func (secret) String() string {
	return "[secret]"
}

var NilSession = Session{}

var ErrNoCurrentUser = errors.New("no current user")

// key value
type kv map[string]interface{}

type resp struct {
	code int
	// FIXME: rename to body
	msg     kv
	session *Session
}

type Server struct {
	sessions *sessions.S
	api      *api.API
	cryptor  *db.Cryptor
}

func New(dataroot, ocrserver string, password []byte) (*mux.Router, error) {

	r := mux.NewRouter().StrictSlash(true)

	km := keysmanager.NewKeysManager(filepath.Join(dataroot, "appsecrets"))

	api := api.NewAPI(dataroot, ocrserver)

	// the datafoot folder doesn't exists, start from scratch
	if _, err := os.Stat(dataroot); os.IsNotExist(err) {
		log.Println("Initiating fresh api...")
		if err := os.Mkdir(dataroot, 0700); err != nil {
			return nil, fmt.Errorf("mkdir %q: %s", dataroot, err)
		}

		if err := km.SignUp(password); err != nil {
			return nil, fmt.Errorf("km.SignUp: %s", err)
		}

		if err := api.Initialize(); err != nil {
			return nil, fmt.Errorf("api.Initialize: %s", err)
		}

	} else {
		log.Printf("Resuming from filesystem...")
		if err := km.Login(password); err != nil {
			return nil, err
		}
		if err := api.Resume(); err != nil {
			return nil, fmt.Errorf("api.Resume: %s", err)
		}
	}

	keys, err := km.LoadKeys()
	if err != nil {
		return nil, fmt.Errorf("loading app keys: %s", err)
	}

	cryptor, err := db.NewCryptor(keys.Encryption, keys.MAC)
	if err != nil {
		return nil, fmt.Errorf("creating cryptor: %s", err)
	}

	sessions, err := sessions.NewS(&sessions.Config{
		// FIXME: keysmanager should be generic (here we are using keys.MAC,
		// but we just need a secure key)
		// SECURITY ISSUE (re use of keys)
		Key: keys.MAC,
	})
	if err != nil {
		return nil, fmt.Errorf("creating sessions.S: %s", err)
	}

	s := &Server{
		sessions: sessions,
		api:      api,
		cryptor:  cryptor,
	}

	rapi := r.PathPrefix("/api").Subrouter()

	rapi.HandleFunc("/", s.h(func(r *http.Request) *resp {
		return &resp{
			code: http.StatusOK,
			msg: kv{
				"kind":        "not implemented",
				"description": "it will list all the different endpoints",
			},
		}
	}))

	//
	// API routes
	//

	post := rapi.Methods(http.MethodPost).Subrouter()

	post.HandleFunc("/login", s.h(s.login))
	post.HandleFunc("/signup", s.h(s.signup))
	post.HandleFunc("/logout", s.h(s.logout))

	post.HandleFunc("/payments/add-manual", s.h(s.addManualPayment))
	rapi.HandleFunc("/payments/list", s.h(s.listPayments))
	rapi.HandleFunc("/payments/scan", s.h(s.scan))

	// make sure this stays at the bottom of the function
	rapi.PathPrefix("/").HandlerFunc(s.h(func(r *http.Request) *resp {
		return &resp{
			code: http.StatusBadRequest,
			msg: kv{
				"kind": "bad request",
				"msg":  "request didn't match any known route",
				"help": []string{
					"make sure you are using the right method (POST instead of GET for example)",
				},
			},
		}
	}))

	html := r.Methods(http.MethodGet).Subrouter()

	html.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./pwa/css"))))
	html.PathPrefix("/icons/").Handler(http.StripPrefix("/icons/", http.FileServer(http.Dir("./pwa/icons"))))

	html.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pwa/manifest.json")
	})

	html.HandleFunc("/js/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Service-Worker-Allowed", "/")
		http.ServeFile(w, r, "./pwa/js/sw.js")
	})
	html.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./pwa/js"))))

	html.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusPermanentRedirect)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// FIXME: implement warning system
			log.Printf("!! warning !! serving %q GET request with html", r.URL.Path)
		}
		http.ServeFile(w, r, "pwa/index.html")
	})

	r.Use(logger)

	return r, nil
}

// h transforms an handlers.handler func to http.HandlerFunc
func (s *Server) h(h func(*http.Request) *resp) http.HandlerFunc {
	handlerName := getFuncName(h)

	return func(w http.ResponseWriter, r *http.Request) {
		resp := h(r)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		if resp.msg == nil {
			resp.code = http.StatusInternalServerError
			// FIXME: implement warning system
			log.Printf("[err] Server handler %s: resp.msg == nil", handlerName)
			resp.msg = kv{
				"kind": "internal error",
				"msg":  "no msg from handler",
			}
		} else if _, ok := resp.msg["kind"]; !ok {
			resp.code = http.StatusInternalServerError
			// FIXME: implement warning system
			log.Printf("[err] Handler %s: no \"kind\" key in resp.msg", handlerName)
			resp.msg = kv{
				"kind": "internal error",
				"msg":  "Handler response was invalid",
			}
		} else {
			if resp.session != nil {
				if err := s.sessions.Save(w, resp.session); err != nil {
					log.Printf("[err] saving session: %s", err)
					resp.code = http.StatusInternalServerError
					resp.msg = kv{
						"kind": "internal error",
						"msg":  "errored saving session cookie",
					}
				}
			} else if resp.session == &NilSession {
				s.sessions.Remove(w)
			}
		}

		if resp.msg["kind"] == "internal error" {
			// FIXME: better error reporting
			log.Printf("!! warning !! internal error: %v %v", r, resp)
		}

		w.WriteHeader(resp.code)
		if err := encoder.Encode(resp.msg); err != nil {
			log.Printf("[err] encoding json object in %s: %s", handlerName, err)
		}
	}
}

func (s *Server) getCurrentUser(r *http.Request) (*db.User, error) {
	session := &Session{}
	err := s.sessions.Load(r, session)
	if errors.Is(err, sessions.ErrInvalidSignature) {
		log.Println("!! Warning !! potential attack on cookie signature")
		return nil, err
	} else if errors.Is(err, sessions.ErrNoSession) {
		return nil, ErrNoCurrentUser
	} else if err != nil {
		return nil, err
	}

	if session.ID == 0 || session.Email == "" || len(session.Password) == 0 {
		log.Println("!! warning !! internal error or potential attack on session")
		log.Println("!! warning !! current session:", session)
		return nil, errors.New("missing fields from session")
	}

	// breaking abstraction here, but I don't know what's better...
	// FIXME: this clearly isn't the right way
	user := db.NewUser(session.ID, session.Email, filepath.Join(s.api.Usersdir, strconv.Itoa(session.ID)))

	password, err := s.cryptor.Decrypt(session.Password)
	if err != nil {
		log.Printf("!! Warning !! decrypting password from session")
		return nil, err
	}

	user.Login(password)

	return user, nil
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%q %q", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
