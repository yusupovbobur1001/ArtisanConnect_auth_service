package postgres

import (
	pb "auth_service/genproto/auth"
	"database/sql"
	"fmt"
)

type UserRepo struct {
	Db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{Db: db}
}

func (u *UserRepo) Register(req *pb.User) error {
	query := `
		insert into users(
				user_name,
				email,
				password,
				full_name,
				user_type,
				bio
		) values($1, $2, $3, $4, $5, $6)`

	_, err := u.Db.Exec(query, req.UserName, req.Email, req.Password,
		req.FullName, req.UserType, req.Bio)

	if err != nil {
		return fmt.Errorf("register metod: %v", err)
	}

	return nil
}


