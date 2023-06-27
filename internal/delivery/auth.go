package delivery

import (
	"errors"
	"forum/internal/model"
	"forum/internal/service"
	"log"
	"net/http"
	"time"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		log.Println("Sign Up: Wrong URL Path")
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.tmpl.ExecuteTemplate(w, "sign_up.html", nil); err != nil {
			log.Printf("Sign Up: Execute: %v", err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("Sign Up: Parse Form: %v", err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		email, ok := r.Form["email"]
		if !ok {
			log.Printf("Sign Up: Parse Form: email field not found")
			h.errorPage(w, http.StatusBadRequest, "email field not found")
			return
		}

		username, ok := r.Form["username"]
		if !ok {
			log.Printf("Sign Up: Parse Form: username field not found")
			h.errorPage(w, http.StatusBadRequest, "username field not found")
			return
		}

		password, ok := r.Form["password"]
		if !ok {
			log.Printf("Sign Up: Parse Form: password field not found")
			h.errorPage(w, http.StatusBadRequest, "password field not found")
			return
		}

		confirmPassword, ok := r.Form["confirm-password"]
		if !ok {
			log.Printf("Sign Up: Parse Form: confirm-password field not found")
			h.errorPage(w, http.StatusBadRequest, "verifyPassword field not found")
			return
		}

		newUser := model.User{
			Email:           email[0],
			Username:        username[0],
			Password:        password[0],
			ConfirmPassword: confirmPassword[0],
		}

		if err := h.Service.Auth.CreateUser(newUser); err != nil {
			log.Printf("Sign Up: Create User: %v", err)
			if errors.Is(err, service.ErrInvalidEmail) ||
				errors.Is(err, service.ErrInvalidUsernameChar) ||
				errors.Is(err, service.ErrConfirmPassword) ||
				errors.Is(err, service.ErrInvalidUsernameLen) ||
				errors.Is(err, service.ErrUserExist) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := h.Service.Auth.GenerateToken(username[0], password[0])
		if err != nil {
			log.Printf("Sign In: Generate Token: %v", err)
			if errors.Is(err, service.ErrUserNotFound) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   user.Token,
			Expires: user.ExpirationTime,
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		log.Println("Sign Up: Method not allowed")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		log.Println("Sign In: Wrong URL Path")
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	switch r.Method {
	case http.MethodGet:
		if err := h.tmpl.ExecuteTemplate(w, "sign_in.html", nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("Sign In: Parse Form: %v", err)
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		username, ok := r.Form["username"]
		if !ok {
			log.Printf("Sign In: Parse Form: email field not found")
			h.errorPage(w, http.StatusBadRequest, "username field not found")
			return
		}

		password, ok := r.Form["password"]
		if !ok {
			log.Printf("Sign In: Parse Form: password field not found")
			h.errorPage(w, http.StatusBadRequest, "password field not found")
			return
		}

		user, err := h.Service.Auth.GenerateToken(username[0], password[0])
		if err != nil {
			log.Printf("Sign In: Generate Token: %v", err)
			if errors.Is(err, service.ErrUserNotFound) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   user.Token,
			Expires: user.ExpirationTime,
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		log.Println("Sign In: Method not allowed")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" {
		log.Println("Logout: Wrong URL Path")
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodGet {
		log.Println("Logout: Method not allowed")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Logout: Get cookie: %v", err)
		if errors.Is(err, http.ErrNoCookie) {
			h.errorPage(w, http.StatusUnauthorized, err.Error())
			return
		}
		h.errorPage(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Service.Auth.DeleteToken(c.Value); err != nil {
		log.Printf("Logout: Delete Token: %v", err)
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Time{},
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
