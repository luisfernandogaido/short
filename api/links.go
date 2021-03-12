package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/luisfernandogaido/short/model"
)

func links(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var body struct {
			Destination string `json:"destination"`
		}
		if err := readJson(r, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		u, err := model.UserToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		link, err := model.LinkCreate(body.Destination, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := struct {
			Link        string `json:"link"`
			Destination string `json:"destination"`
		}{domain + "/" + link.Hash, link.Destination}
		printJson(w, res)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		hash := strings.Replace(r.URL.Path, "/", "", 1)
		link, err := model.LinkGet(hash)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, link.Destination, http.StatusMovedPermanently)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
