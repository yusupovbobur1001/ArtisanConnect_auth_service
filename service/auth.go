package service

import (
	pb "auth_service/genproto/auth"
	"auth_service/storage/postgres"
	"context"
	"database/sql"
	"fmt"
)

type Handler struct {
	Auth *postgres.UserRepo
	pb.UnimplementedAuthServer
}

func NewHadler(db *sql.DB) *Handler {
	return &Handler{Auth: postgres.NewUserRepo(db)}
}

func (h *Handler) UpdateProfile(ctx context.Context, req *pb.UserUpdate) (*pb.GetProfile, error) {
	if req.UserType != "verdor" || req.UserType != "receiver" {
		return nil, fmt.Errorf("UserType biz hohlaganday emas!")
	}
	userP, err := h.Auth.UpdateUser(req)
}
