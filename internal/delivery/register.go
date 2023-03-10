package delivery

import (
	"database/sql"
	"forum/internal/service"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmp, err := template.ParseFiles("./ui/html/register.html")
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, "email")
	case http.MethodPost:
		tmp, err := template.ParseFiles("./ui/html/login.html")
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		crpassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		password = string(crpassword)
		database, _ := sql.Open("sqlite3", "./forum.db")
		defer database.Close()
		DB, _ := database.Prepare(`Insert into users(email, username, password) values(?, ?, ?)`)
		rows, _ := database.Query("select * from users where email ='" + email + "' or username ='" + username + "'")
		var id int
		var username2 string
		var email2 string
		var password2 string
		for rows.Next() {
			rows.Scan(&id, &email2, &username2, &password2)
		}
		if username2 != "" && email2 != "" && password2 != "" || !service.ValidMailAddress(email) {
			tmp, err := template.ParseFiles("./ui/html/register.html")
			if err != nil {
				Errors(w, http.StatusInternalServerError)
				return
			}
			tmp.Execute(w, "change your username")
		} else {
			DB.Exec(email, username, password)
			tmp.Execute(w, nil)
		}
	}
}
