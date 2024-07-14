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
	userP, err := h.Auth.UpdateUser(req)
	if err != nil {
		return nil, err
	}
	return userP, nil
}

func (h *Handler) DeleteProfile(ctx context.Context, req *pb.Id) (*pb.Message, error) {
	resp, err := h.Auth.DeleteUser(req)
	if err != nil {
		return &pb.Message{Message: "The user was not successfully deleted"}, err
	}
	return resp, nil
}

func (h *Handler) GetByIdProfile(ctx context.Context, req *pb.Id) (*pb.UsersInfo, error) {
	h.Auth.GetByIdUser(req)
}
