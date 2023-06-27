package delivery

import (
	"context"
	"errors"
	"forum/internal/model"
	"net/http"
	"time"
)

const ctxKeyUser ctxKey = iota

type ctxKey int8

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var err error
		c, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, model.User{})))
				return
			}
			h.errorPage(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err = h.Service.ParseToken(c.Value)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, model.User{})))
			return
		}
		if user.ExpirationTime.Before(time.Now()) {
			if err := h.Service.DeleteToken(c.Value); err != nil {
				h.errorPage(w, http.StatusInternalServerError, err.Error())
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, model.User{})))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}
