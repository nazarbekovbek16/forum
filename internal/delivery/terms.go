package delivery

import (
	"html/template"
	"net/http"
)

func Terms(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/terms" {
		Errors(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed)
		return
	}
	tmp, err := template.ParseFiles("./ui/html/terms.html")
	if err != nil {
		Errors(w, http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}
