package delivery

import (
	"forum/internal/service"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl    *template.Template
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		tmpl:    template.Must(template.ParseGlob("web/template/*.html")),
		Service: service,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.userIdentity(h.homePage))

	mux.HandleFunc("/auth/signup", h.signUp)
	mux.HandleFunc("/auth/signin", h.signIn)
	mux.HandleFunc("/auth/logout", h.logout)

	mux.HandleFunc("/post/", h.userIdentity(h.postPage))
	mux.HandleFunc("/post/create", h.userIdentity(h.createPost))
	mux.HandleFunc("/post/like/", h.userIdentity(h.likePost))
	mux.HandleFunc("/post/dislike/", h.userIdentity(h.dislikePost))

	mux.HandleFunc("/comment/like/", h.userIdentity(h.likeComment))
	mux.HandleFunc("/comment/dislike/", h.userIdentity(h.dislikeComment))

	mux.HandleFunc("/profile/", h.userIdentity(h.userProfile))

	mux.Handle("/static/css/", http.StripPrefix("/static/css", http.FileServer(http.Dir("./web/static/css"))))
	mux.Handle("/static/img/", http.StripPrefix("/static/img", http.FileServer(http.Dir("./web/static/img"))))
}
