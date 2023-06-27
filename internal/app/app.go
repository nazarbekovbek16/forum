package app

import (
	"forum/internal/config"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"log"
)

func Run(cfgFilePath string) error {
	cfg := config.NewConfig(cfgFilePath)

	db, err := repository.InitDB(cfg)
	if err != nil {
		return err
	}

	defer db.Close()

	if err := repository.CreateTables(db); err != nil {
		return err
	}

	repository := repository.NewRepository(db, cfg)
	service := service.NewService(repository)
	handler := delivery.NewHandler(service)

	server := server.NewServer(cfg, handler)
	log.Printf("Starting server at port %v\nhttp://localhost:%v\n", cfg.Http.Addr, cfg.Http.Addr)

	return server.Srv.ListenAndServe()
}
