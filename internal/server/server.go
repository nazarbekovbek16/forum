package server

import (
	"forum/internal/delivery"
	"forum/internal/repository"
)

func Server() {
	repository.CreateTable()
	delivery.Handlers()
}
