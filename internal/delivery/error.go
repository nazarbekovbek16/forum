package delivery

import (
	"net/http"
	"strings"
)

func (h *Handler) errorPage(w http.ResponseWriter, code int, errorText string) {
	w.WriteHeader(code)
	data := struct {
		Status  int
		Message string
		ErrText string
	}{
		Status:  code,
		Message: http.StatusText(code),
		ErrText: errorText,
	}
	if data.Status != http.StatusInternalServerError {
		temp := strings.Split(errorText, ":")
		data.ErrText = temp[len(temp)-1]
	}
	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// fmt.Fprintf(w, "%d - %s\n", data.Status, data.Message)
	}
}
