package handler

import (
	"auth_service/pkg/logger"
	"auth_service/storage/postgres"
	"database/sql"
	"log"
	"log/slog"
)

type Handler struct {
	Auth   postgres.UserRepo
	Logger *slog.Logger
}

func NewHadler(db *sql.DB) *Handler {
	l, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	return &Handler{
		Logger: l,
		Auth: *postgres.NewUserRepo(db),
	}
}