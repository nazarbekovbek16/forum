package delivery

import (
	"database/sql"
	"net/http"
	"strconv"
)

func Like(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			Errors(w, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	like := r.FormValue("like")
	dislike := r.FormValue("dislike")

	sessionToken := cook.Value
	Id := r.URL.Query().Get("id")
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	sessions := database.QueryRow("select user from sessions where session ='" + sessionToken + "'")
	var user string
	sessions.Scan(&user)

	if dislike == "dislike" {
		row := database.QueryRow("select * from likes where postID =" + Id + " and owner='" + user + "'")
		var id int
		var postID int
		var owner string
		var likes int
		var dislikes int
		row.Scan(&id, &postID, &owner, &likes, &dislikes)
		if owner == "" {
			dislikes = 1
			DB, _ := database.Prepare(`Insert into likes(postID, owner, like, dislike) values(?, ?, ?, ?)`)
			DB.Exec(Id, user, likes, dislikes)
		} else {
			DB, _ := database.Prepare("update likes set dislike=1, like=0 where id=" + strconv.Itoa(id) + "")
			DB.Exec()
		}
	} else if like == "like" {
		row := database.QueryRow("select * from likes where postID =" + Id + " and owner='" + user + "'")
		var id int
		var postID int
		var owner string
		var likes int
		var dislikes int
		row.Scan(&id, &postID, &owner, &likes, &dislikes)
		if owner == "" {
			likes = 1
			DB, _ := database.Prepare(`Insert into likes(postID, owner, like, dislike) values(?, ?, ?, ?)`)
			DB.Exec(Id, user, likes, dislikes)
		} else {
			DB, _ := database.Prepare("update likes set dislike=0, like=1 where id=" + strconv.Itoa(id) + "")
			DB.Exec()
		}
	}

	http.Redirect(w, r, "/read/?id="+Id, http.StatusFound)
}
