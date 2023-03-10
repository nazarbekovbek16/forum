package delivery

import (
	"database/sql"
	"forum/internal/model"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Read(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		Nick := false
		cook, err := r.Cookie(COOKIE_NAME)
		if err == nil {
			Nick = true

			sess := cook.Value

			DB, _ := sql.Open("sqlite3", "./forum.db")
			defer DB.Close()
			rows, _ := DB.Query("select * from sessions where session ='" + sess + "'")
			var id int
			var user string
			var session string
			for rows.Next() {
				rows.Scan(&id, &user, &session)
			}

			if user == "" {
				Nick = false
			}
		}
		tmp, err := template.ParseFiles("./ui/html/read.html")
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}

		Id := r.URL.Query().Get("id")
		database, _ := sql.Open("sqlite3", "./forum.db")
		defer database.Close()
		DBlikes := database.QueryRow("select count(*) from likes where postID=" + Id + " and like=1")
		DBdislikes := database.QueryRow("select count(*) from likes where postID=" + Id + " and dislike=1")
		var like int
		var dislike int
		DBlikes.Scan(&like)
		DBdislikes.Scan(&dislike)
		DB, _ := database.Prepare("update post set likes=" + strconv.Itoa(like) + ", dislikes=" + strconv.Itoa(dislike) + " where id='" + Id + "'")
		DB.Exec()
		rows := database.QueryRow("select * from post where id =" + Id)
		var id int
		var owner string
		var title string
		var content string
		var category string
		var image string
		var likes int
		var dislikes int

		rows.Scan(&id, &owner, &title, &content, &category, &image, &likes, &dislikes)
		if id == 0 {
			Errors(w, http.StatusBadRequest)
			return
		}
		items := model.PostItem{
			ID:       id,
			Owner:    owner,
			Title:    title,
			Content:  content,
			Types:    category,
			Image:    image,
			Likes:    likes,
			Dislikes: dislikes,
		}

		itemsComment := []model.CommentsItem{}
		comments, _ := database.Query("SELECT * FROM comment where postID='" + Id + "'")
		var idComments int
		var postId int
		var ownerComments string
		var contentComments string

		for comments.Next() {
			comments.Scan(&idComments, &postId, &ownerComments, &contentComments)
			DBlikesComment := database.QueryRow("select count(*) from likescomment where commentID=" + strconv.Itoa(idComments) + " and like=1")
			DBdislikesComment := database.QueryRow("select count(*) from likescomment where commentID=" + strconv.Itoa(idComments) + " and dislike=1")
			var likeComment int
			var dislikeComment int
			DBlikesComment.Scan(&likeComment)
			DBdislikesComment.Scan(&dislikeComment)
			itemComment := model.CommentsItem{
				ID:       idComments,
				PostId:   postId,
				Owner:    ownerComments,
				Comment:  contentComments,
				Likes:    likeComment,
				Dislikes: dislikeComment,
			}
			itemsComment = append(itemsComment, itemComment)
		}
		res := model.Read{
			Posts:       items,
			CurrentUser: Nick,
			Comments:    itemsComment,
		}
		tmp.Execute(w, res)
	case http.MethodPost:
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

		Id := r.URL.Query().Get("id")
		comment := r.FormValue("comment")
		checkComm := strings.TrimSpace(comment)
		if comment != "" && user != "" && checkComm != "" {
			DB, _ := database.Prepare(`Insert into comment(postID, owner, content) values(?, ?, ?)`)
			DB.Exec(Id, user, comment)
			http.Redirect(w, r, "/read/?id="+Id, http.StatusFound)
		} else {
			http.Redirect(w, r, "/read/?id="+Id, http.StatusFound)
		}

	}
}
