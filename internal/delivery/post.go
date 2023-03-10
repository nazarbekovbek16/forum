package delivery

import (
	"html/template"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmp, err := template.ParseFiles("./ui/html/post.html")
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, nil)
	}
}
