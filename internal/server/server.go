package server

import (
	"forum/internal/config"
	"forum/internal/delivery"
	"net/http"
	"time"
)

type Server struct {
	Srv *http.Server
}

func NewServer(cfg *config.Config, handler *delivery.Handler) *Server {
	mux := http.NewServeMux()
	handler.InitRoutes(mux)
	return &Server{
		Srv: &http.Server{
			Addr:           ":" + cfg.Http.Addr,
			Handler:        mux,
			ReadTimeout:    time.Duration(time.Duration(cfg.Http.ReadTimeout).Seconds()),
			WriteTimeout:   time.Duration(time.Duration(cfg.Http.WriteTimeout).Seconds()),
			MaxHeaderBytes: cfg.Http.MaxHeaderByte,
		},
	}
}
