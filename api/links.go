package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/luisfernandogaido/short/model"
)

func links(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var body struct {
			Destination string `json:"destination"`
			Hash        string `json:"hash"`
			TtlDays     int    `json:"ttl_days"`
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
		domainFound := false
		for _, ad := range authorizedDomains {
			if strings.Contains(body.Destination, ad) {
				domainFound = true
				break
			}
		}
		if !domainFound {
			http.Error(w, "destino não autorizado", http.StatusForbidden)
			return
		}
		link, err := model.LinkCreate(body.Destination, body.Hash, body.TtlDays, u)
		if err != nil {
			if errors.Is(err, model.ErrDuplicated) {
				http.Error(w, "hash em uso", http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := struct {
			Link        string    `json:"link"`
			Destination string    `json:"destination"`
			ExpiresAt   time.Time `json:"expires_at"`
		}{domain + "/" + link.Hash, link.Destination, link.ExpiresAt.Local()}
		printJson(w, res)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	defer func() {
		go loga(r)
	}()
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
