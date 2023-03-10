package main

import (
	"forum/internal/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	server.Server()
}
