package api

import (
	"net/http"
	"strings"

	"github.com/luisfernandogaido/short/model"
)

func users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, err := model.Users()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		printJson(w, users)
	case "POST":
		var user model.User
		if err := readJson(r, &user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := user.Save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		printJson(w, user)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func user(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := strings.Replace(r.URL.Path, "/users/", "", 1)
		u, err := model.NewUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		printJson(w, u)
	case "PUT":
		id := strings.Replace(r.URL.Path, "/users/", "", 1)
		u, err := model.NewUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var body model.User
		if err := readJson(r, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.Name = body.Name
		if err := u.Save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		printJson(w, u)
	case "DELETE":
		id := strings.Replace(r.URL.Path, "/users/", "", 1)
		u, err := model.NewUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := u.Delete(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func userRegenerateToken(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		id := strings.Replace(r.URL.Path, "/users/regen-token/", "", 1)
		u, err := model.NewUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := u.RegenerateToken(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		printJson(w, u)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
