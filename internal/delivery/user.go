package delivery

import (
	"errors"
	"forum/internal/model"
	"forum/internal/service"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) userProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(model.User)
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	userPage, err := h.Service.User.GetUserByUsername(username)
	if err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	if r.Method != http.MethodGet {
		log.Println(err)
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	posts, err := h.Service.User.GetPostByUsername(userPage.Username, r.URL.Query())
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrInvalidQuery) {
			h.errorPage(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	info := model.Info{
		User:        user,
		ProfileUser: userPage,
		Posts:       posts,
	}

	if err := h.tmpl.ExecuteTemplate(w, "profile.html", info); err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
