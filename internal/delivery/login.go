package delivery

import (
	"database/sql"
	"forum/internal/service"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var InMemorySession *service.Session

const COOKIE_NAME = "sessionId"

func Login(w http.ResponseWriter, r *http.Request) {
	InMemorySession = service.NewSession()
	switch r.Method {
	case http.MethodGet:
		tmp, err := template.ParseFiles("./ui/html/login.html")
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, nil)
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.Form.Get("password")

		database, _ := sql.Open("sqlite3", "./forum.db")
		defer database.Close()
		rows, _ := database.Query("select * from users where email like '" + email + "'")
		var id int
		var username2 string
		var email2 string
		var password2 string
		for rows.Next() {
			rows.Scan(&id, &email2, &username2, &password2)
		}
		err := bcrypt.CompareHashAndPassword([]byte(password2), []byte(password))
		//////////// SESSION //////////////////
		if username2 != "" && email2 != "" && password2 != "" && err == nil {

			sessionId := InMemorySession.Init(username2)
			Nick := InMemorySession.Get(sessionId)
			rows, _ := database.Query("select * from sessions where user ='" + Nick + "'")
			var id int
			var user string
			var session string
			for rows.Next() {
				rows.Scan(&id, &user, &session)
			}
			if Nick == user {
				DB, _ := database.Prepare("update sessions set session=? where user=?")
				DB.Exec(sessionId, Nick)
			} else {
				DB, _ := database.Prepare(`Insert into sessions(user,session) values(?,?)`)
				DB.Exec(Nick, sessionId)
			}
			cookie := &http.Cookie{
				Name:    COOKIE_NAME,
				Value:   sessionId,
				Expires: time.Now().Add(5 * time.Minute),
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}
