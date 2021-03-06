package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/luisfernandogaido/short/model"
)

var (
	root              string
	domain            string
	authorizedDomains []string
)

func Serve(addr string, mongoURI string, tokenRoot string, redisURI string, dom string, ads []string) error {
	root = tokenRoot
	domain = dom
	authorizedDomains = ads

	if err := model.Ini(mongoURI, redisURI); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	go model.LinkPurge()

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/users", autRoot(users))
	http.HandleFunc("/users/", autRoot(user))
	http.HandleFunc("/users/regen-token/", autRoot(userRegenerateToken))
	http.HandleFunc("/links", aut(links))
	http.HandleFunc("/", redirect)

	return http.ListenAndServe(addr, nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func autRoot(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if t != root {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		f(w, r)
		go loga(r)
	}
}

func aut(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if err := model.Aut(t); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		f(w, r)
		go loga(r)
	}
}

func loga(r *http.Request) {
	t := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}
	err := model.AcessoLoga(model.Acesso{
		Ip:        ip,
		Token:     t,
		Data:      time.Now(),
		Path:      r.URL.Path,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		log.Println(err)
	}
}

func printJson(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("printjson: %w", err)
	}
	return nil
}

func readJson(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return fmt.Errorf("readjson: %w", err)
	}
	return nil
}
