package delivery

import (
	"database/sql"
	"forum/internal/service"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c, err := r.Cookie("sessionId")
		if err != nil {
			if err == http.ErrNoCookie {

				Errors(w, http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value

		database, _ := sql.Open("sqlite3", "./forum.db")
		defer database.Close()
		rows := database.QueryRow("select * from sessions where session ='" + sessionToken + "'")
		var id int
		var user string
		var session string
		rows.Scan(&id, &user, &session)

		InMemorySession = service.NewSession()

		database.Exec("delete from sessions where id =$1", id)

		http.SetCookie(w, &http.Cookie{
			Name:    COOKIE_NAME,
			Value:   "",
			Expires: time.Now(),
		})

		http.Redirect(w, r, "/", http.StatusFound)
	case http.MethodPost:
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
