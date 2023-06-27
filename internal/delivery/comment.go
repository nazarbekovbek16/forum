package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/model"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(model.User)
	if r.Method != http.MethodPost {
		log.Println("method mot allowed like comment")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/like/"))
	if err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	comment, err := h.Service.Commentary.GetCommentaryById(id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.Service.VoteComment.LikeCommentary(id, user.Username); err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostID), http.StatusSeeOther)
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(model.User)
	if r.Method != http.MethodPost {
		log.Println("method not allowed dislike comment")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/dislike/"))
	if err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	comment, err := h.Service.Commentary.GetCommentaryById(id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.Service.VoteComment.DislikeCommentary(id, user.Username); err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostID), http.StatusSeeOther)
}
