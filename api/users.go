package api

import (
	"fmt"
	"net/http"

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
		fmt.Fprintln(w, "POST não implementado ainda")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

}
