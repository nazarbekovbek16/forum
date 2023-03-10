package delivery

import (
	"database/sql"
	"net/http"
	"strconv"
)

func LikeComment(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			Errors(w, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	like := r.FormValue("likeC")
	dislike := r.FormValue("dislikeC")

	sessionToken := cook.Value
	Id := r.URL.Query().Get("id")
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	findID := database.QueryRow("select postID from comment where id='" + Id + "'")
	var redirect int
	findID.Scan(&redirect)
	sessions := database.QueryRow("select user from sessions where session ='" + sessionToken + "'")
	var user string
	sessions.Scan(&user)

	if dislike == "dislike" {
		row := database.QueryRow("select * from likescomment where commentID =" + Id + " and owner='" + user + "'")
		var id int
		var commentID int
		var owner string
		var likes int
		var dislikes int
		row.Scan(&id, &commentID, &owner, &likes, &dislikes)
		if owner == "" {
			dislikes = 1
			DB, _ := database.Prepare(`Insert into likescomment(commentID, owner, like, dislike) values(?, ?, ?, ?)`)
			DB.Exec(Id, user, likes, dislikes)
		} else {
			DB, _ := database.Prepare("update likescomment set dislike=1, like=0 where id=" + strconv.Itoa(id) + "")
			DB.Exec()
		}
	} else if like == "like" {
		row := database.QueryRow("select * from likescomment where commentID =" + Id + " and owner='" + user + "'")
		var id int
		var commentID int
		var owner string
		var likes int
		var dislikes int
		row.Scan(&id, &commentID, &owner, &likes, &dislikes)
		if owner == "" {
			likes = 1
			DB, _ := database.Prepare(`Insert into likescomment(commentID, owner, like, dislike) values(?, ?, ?, ?)`)
			DB.Exec(Id, user, likes, dislikes)
		} else {
			DB, _ := database.Prepare("update likescomment set dislike=0, like=1 where id=" + strconv.Itoa(id) + "")
			DB.Exec()
		}
	}

	http.Redirect(w, r, "/read/?id="+strconv.Itoa(redirect), http.StatusFound)
}
