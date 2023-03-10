package delivery

import (
	"html/template"
	"net/http"
)

func Errors(w http.ResponseWriter, code int) {
	tml, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		http.Error(w, "Internal Error", 500)
		return
	}
	w.WriteHeader(code)
	err = tml.Execute(w, http.StatusText(code))
	if err != nil {
		http.Error(w, "Internal Error", 500)
		return
	}
}
