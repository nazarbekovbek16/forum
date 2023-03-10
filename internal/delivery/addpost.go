package delivery

import (
	"database/sql"
	"html/template"
	"net/http"
)

func Addpost(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			Errors(w, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := cook.Value
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	rows, _ := database.Query("select * from sessions where session ='" + sessionToken + "'")
	var id int
	var user string
	var session string
	for rows.Next() {
		rows.Scan(&id, &user, &session)
	}
	if user != "" {
		title := r.FormValue("newtitle")
		content := r.FormValue("newcontent")
		typeTemp := r.Form["newcategory"]
		types := ""
		for _, categ := range typeTemp {
			types += categ + " "
		}

		image := r.FormValue("newimage")
		if title != "" && content != "" {
			database, _ := sql.Open("sqlite3", "./forum.db")
			defer database.Close()
			DB, _ := database.Prepare(`INSERT INTO post (owner,title,content,type,image) values (?,?,?,?,?)`)
			DB.Exec(user, title, content, types, image)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/post/", http.StatusFound)
		}

	} else {
		tmp, err := template.ParseFiles("./ui/html/login.html")
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, nil)
	}
}
