package delivery

import (
	"log"
	"net/http"
	"os"
)

func Handlers() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/register", Register)
	mux.HandleFunc("/post/", Post)
	mux.HandleFunc("/addpost", Addpost)
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/terms", Terms)
	mux.HandleFunc("/read/", Read)
	mux.HandleFunc("/like/", Like)
	mux.HandleFunc("/likecomment/", LikeComment)
	fileServer := http.FileServer(http.Dir("./ui/style"))
	mux.Handle("/style", http.NotFoundHandler())
	mux.Handle("/style/", http.StripPrefix("/style", fileServer))
	log.Println("Запуск веб-сервера на http://127.0.0.1:8000")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
