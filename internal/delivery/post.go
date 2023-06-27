package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) postPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(model.User)

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		log.Println("Post not found")
		h.errorPage(w, http.StatusBadRequest, "post not found")
		return
	}

	post, err := h.Service.Post.GetPostByID(id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		comments, err := h.Service.GetCommentariesByPostID(post.ID)
		if err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		postLikes, err := h.Service.GetPostLikes(post.ID)
		if err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		postDisikes, err := h.Service.GetPostDislikes(post.ID)
		if err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentsLikes, err := h.Service.GetCommentaryLikes(post.ID)
		if err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentsDislikes, err := h.Service.GetCommentaryDislikes(post.ID)
		if err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		info := model.Info{
			Post:             post,
			PostLikes:        postLikes,
			PostDislikes:     postDisikes,
			User:             user,
			Commentaries:     comments,
			CommentsLikes:    commentsLikes,
			CommentsDislikes: commentsDislikes,
		}
		if err := h.tmpl.ExecuteTemplate(w, "post.html", info); err != nil {
			log.Printf("Post page: Executing %v", err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		if user == (model.User{}) {
			h.errorPage(w, http.StatusUnauthorized, "cant post comment")
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("Post page: Parse Form: %v", err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		comment, ok := r.Form["comment"]
		if !ok {
			log.Printf("Post page: Parse Form: comment field not found")
			h.errorPage(w, http.StatusBadRequest, "comment field not foud")
			return
		}

		newComment := model.Commentary{
			PostID:  post.ID,
			Author:  user.Username,
			Content: comment[0],
		}

		if err := h.Service.Commentary.CreateCommentary(newComment); err != nil {
			log.Println(err)
			if errors.Is(err, service.ErrInvalidComment) ||
				errors.Is(err, service.ErrCommentLen) || errors.Is(err, service.ErrInvalidCommentChar) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		log.Println("Method not allowed post")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(model.User)
	if user == (model.User{}) {
		h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	if r.URL.Path != "/post/create" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	switch r.Method {
	case http.MethodGet:
		info := model.Info{
			User: user,
		}

		if err := h.tmpl.ExecuteTemplate(w, "create_post.html", info); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		title, ok := r.Form["title"]
		if !ok {
			log.Println("title field not found")
			h.errorPage(w, http.StatusBadRequest, "title field not found")
			return
		}

		content, ok := r.Form["content"]
		if !ok {
			log.Println("content field not found")
			h.errorPage(w, http.StatusBadRequest, "content field not found")
			return
		}

		category, ok := r.Form["categories"]
		if !ok {
			log.Println("category field not found")
			h.errorPage(w, http.StatusBadRequest, "category field not found")
			return
		}
		post := model.Post{
			Title:    title[0],
			Content:  content[0],
			Author:   user.Username,
			Category: category,
		}

		if err := h.Service.Post.CreatePost(post); err != nil {
			log.Println(err)
			if errors.Is(err, service.ErrInvalidPostContent) || errors.Is(err, service.ErrInvalidPostTitle) || errors.Is(err, service.ErrPostContentLen) || errors.Is(err, service.ErrPostTitleLen) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(model.User)
	if user == (model.User{}) {
		h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/like/"))
	if err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	if r.Method != http.MethodPost {
		log.Println("Method not allowed likepost")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if err := h.Service.VotePost.LikePost(id, user.Username); err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *Handler) dislikePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(model.User)

	if user == (model.User{}) {
		h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/dislike/"))
	if err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	if r.Method != http.MethodPost {
		log.Println("Method not allowed likepost")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if err := h.Service.VotePost.DislikePost(id, user.Username); err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}
