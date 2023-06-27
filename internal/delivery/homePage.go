package delivery

import (
	"errors"
	"forum/internal/model"
	"forum/internal/service"
	"log"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	var posts []model.Post
	var err error
	if r.URL.Path != "/" {
		// log.Printf("home page: wrong url:\n")
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		// log.Printf("home page: get post by filter:\n")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(ctxKeyUser).(model.User)

	if len(r.URL.Query()) == 0 {
		posts, err = h.Service.Post.GetAllPosts()
		if err != nil {
			// log.Println(err)
			log.Printf("home page: check query len: %v \n", err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		posts, err = h.Service.Post.GetAllPostsByFilter(user, r.URL.Query())
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				log.Printf("home page: get post by filter: %v \n", err)
				h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	info := model.Info{
		Posts: posts,
		User:  user,
	}

	if err := h.tmpl.ExecuteTemplate(w, "index.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
