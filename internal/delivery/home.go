package delivery

import (
	"database/sql"
	"forum/internal/model"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/" {
			Errors(w, http.StatusNotFound)
			return
		}
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

		tmp, err := template.ParseFiles("./ui/html/index.html")

		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}

		DB, _ := sql.Open("sqlite3", "./forum.db")
		defer DB.Close()
		category := r.URL.Query().Get("type")

		if category != "" {
			items1 := []model.PostItem{}

			var rows1 *sql.Rows
			if category == "mynews" {
				cook, err := r.Cookie(COOKIE_NAME)
				if err != nil {
					Errors(w, http.StatusUnauthorized)
					return
				}
				sessionToken := cook.Value

				rows, _ := DB.Query("select * from sessions where session ='" + sessionToken + "'")
				var id int
				var user string
				var session string
				for rows.Next() {
					rows.Scan(&id, &user, &session)
				}
				rows1, err = DB.Query("SELECT * FROM post where owner='" + user + "'")
			} else if category == "liked" {
				cook, err := r.Cookie(COOKIE_NAME)
				if err != nil {
					Errors(w, http.StatusUnauthorized)
					return
				}
				sessionToken := cook.Value

				rows, _ := DB.Query("select * from sessions where session ='" + sessionToken + "'")
				var id int
				var user string
				var session string
				for rows.Next() {
					rows.Scan(&id, &user, &session)
				}

				var likedPost []int
				rowsLike, _ := DB.Query("select postID from likes where owner ='" + user + "'")
				for rowsLike.Next() {
					var postID int
					rowsLike.Scan(&postID)
					likedPost = append(likedPost, postID)
				}
				items2 := []model.PostItem{}
				var id1 int
				var owner1 string
				var title1 string
				var content1 string
				var types1 string
				var image1 string
				var likes1 int
				var dislikes1 int
				for i := 0; i < len(likedPost); i++ {
					rowItem := DB.QueryRow("select * from post where id='" + strconv.Itoa(likedPost[i]) + "'")
					rowItem.Scan(&id1, &owner1, &title1, &content1, &types1, &image1, &likes1, &dislikes1)
					item1 := model.PostItem{
						ID:       id1,
						Owner:    owner1,
						Title:    title1,
						Content:  content1,
						Types:    types1,
						Image:    image1,
						Likes:    likes1,
						Dislikes: dislikes1,
					}
					items2 = append(items2, item1)
				}
				res := model.Home{
					Posts:       items2,
					CurrentUser: Nick,
				}
				err = tmp.Execute(w, res)
				return
			} else {
				rows1, err = DB.Query("SELECT * FROM post")
				var id1 int
				var owner1 string
				var title1 string
				var content1 string
				var types1 string
				var image1 string
				var likes1 int
				var dislikes1 int
				for rows1.Next() {
					rows1.Scan(&id1, &owner1, &title1, &content1, &types1, &image1, &likes1, &dislikes1)
					item1 := model.PostItem{
						ID:       id1,
						Owner:    owner1,
						Title:    title1,
						Content:  content1,
						Types:    types1,
						Image:    image1,
						Likes:    likes1,
						Dislikes: dislikes1,
					}
					temp := strings.Split(types1, " ")
					for _, categ := range temp {
						if categ == category {
							items1 = append(items1, item1)
						}
					}

				}
				res := model.Home{
					Posts:       items1,
					CurrentUser: Nick,
				}
				err = tmp.Execute(w, res)
				if err != nil {
					Errors(w, http.StatusInternalServerError)
					return
				}
				return
			}
			var id1 int
			var owner1 string
			var title1 string
			var content1 string
			var types1 string
			var image1 string
			var likes1 int
			var dislikes1 int
			for rows1.Next() {
				rows1.Scan(&id1, &owner1, &title1, &content1, &types1, &image1, &likes1, &dislikes1)
				item1 := model.PostItem{
					ID:       id1,
					Owner:    owner1,
					Title:    title1,
					Content:  content1,
					Types:    types1,
					Image:    image1,
					Likes:    likes1,
					Dislikes: dislikes1,
				}
				items1 = append(items1, item1)

			}
			res := model.Home{
				Posts:       items1,
				CurrentUser: Nick,
			}
			err = tmp.Execute(w, res)
			if err != nil {
				Errors(w, http.StatusInternalServerError)
				return
			}
		} else {
			items := []model.PostItem{}
			rows, _ := DB.Query(`SELECT * FROM post`)
			var id int
			var owner string
			var title string
			var content string
			var types string
			var image string
			var likes int
			var dislikes int
			for rows.Next() {
				rows.Scan(&id, &owner, &title, &content, &types, &image, &likes, &dislikes)
				item := model.PostItem{
					ID:       id,
					Owner:    owner,
					Title:    title,
					Content:  content,
					Types:    types,
					Image:    image,
					Likes:    likes,
					Dislikes: dislikes,
				}
				items = append(items, item)
			}
			res := model.Home{
				Posts:       items,
				CurrentUser: Nick,
			}
			err = tmp.Execute(w, res)
			if err != nil {
				Errors(w, http.StatusInternalServerError)
				return
			}
		}
	case http.MethodPost:
		tmp, err := template.ParseFiles("./ui/html/index.html")
		err = tmp.Execute(w, nil)
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}
	}
}
