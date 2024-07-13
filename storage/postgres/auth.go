package postgres

import (
	pb "auth_service/genproto/auth"
	"auth_service/model"
	"database/sql"
	"fmt"
	"time"
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


func (u *UserRepo) Login(req *pb.UserLogin) (*model.User, error) {
	user := model.User{}
	query := `
		select 
			user_name, full_name, user_type, bio
		from 
			users
		where email=$1 and password=$2`
	
	err := u.Db.QueryRow(query, req.Email, req.Password).Scan(&user.UserName, &user.FullName, &user.UserName, &user.Bio)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (u *UserRepo) UpdateUser(req *pb.UserUpdate, id string) (*pb.GetProfile, error) {
	query := `
		update 
			users
		set
			user_name = $1, full_name = $2, bio = $3, user_type = $4, updated_at = $6
		where id = $5`
	
	_, err := u.Db.Exec(query, req.UserName, req.FullName, req.Bio,  req.UserType, id, time.Now())
	if err != nil {
		return nil, err
	}

	resp, err := u.GetByIdUser(&pb.Id{Id: id})
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (u *UserRepo) DeleteUser(req *pb.Id) (*pb.Message, error) {
	_, err := u.Db.Exec(`
						update 
							users
						set
							deleted_at = $1
						where 
							id = $2`, time.Now(), req.Id)
	
	if err != nil {
		return &pb.Message{Message: "The user was not successfully deleted"}, err
	}
	return &pb.Message{Message: "User successfully deleted"}, nil
}

func (u *UserRepo) GetByIdUser(req *pb.Id) (*pb.UserInfo, error) {
	user := pb.UserInfo{}
	query := `
		serlect 
			user_name, email, password, full_name, user_type, id, created_at, bio
		from 
			uesrs
		where 
			id = $1 and deleted_at is null`
	
	err := u.Db.QueryRow(query, req.Id).Scan(&user.UserName, &user.Email, &user.Password,
						&user.FullName, &user.UserType, &user.Id, &user.CreatedAt, &user.Bio)
	if err != nil {
		return nil, err
	}
	return &user, nil
}