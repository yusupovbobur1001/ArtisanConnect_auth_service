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
				bio,
				created_at
		) values($1, $2, $3, $4, $5, $6, $7)`

	_, err := u.Db.Exec(query, req.UserName, req.Email, req.Password,
		req.FullName, req.UserType, req.Bio, time.Now())

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

func (u *UserRepo) Logout(token *pb.Tokens) error {
	_, err := u.Db.Exec(`
	update 
		reflesh_tokens 
	set 
		deleted_at=$1 
	where 
		token=$2`, time.Now(), token.RefreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) RefreshToken(refresh string) (bool, error) {
	query := `
	select 
		case 
			when token = $1 then true 
		else 
			false 
		end 
	from 
		refresh_tokens 
	where 
		token = $1 and deletad_at is null
	`
	exists := false
	err := u.Db.QueryRow(query, refresh).Scan(&exists)

	return exists, err
}

func (u *UserRepo) UpdateUser(req *pb.UserUpdate) (*pb.GetProfile, error) {
	query := `
		update 
			users
		set
			user_name = $1, 
			full_name = $2, 
			bio = $3, 
			user_type = $4, 
			updated_at = $6
		where 
			id = $5 `

	_, err := u.Db.Exec(query, req.UserName, req.FullName, req.Bio, req.UserType, req.Id, time.Now())
	if err != nil {
		return nil, err
	}

	resp, err := u.GetByIdUser(&pb.Id{Id: req.Id})
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

func (u *UserRepo) GetByIdUser(req *pb.Id) (*pb.GetProfile, error) {
	user := pb.GetProfile{}
	query := `
		select 
			user_name, email, password, full_name, id, updated_at, bio
		from 
			users
		where 
			id = $1 and deleted_at is null`

	err := u.Db.QueryRow(query, req.Id).Scan(&user.UserName, &user.Email, &user.Password,
		&user.FullName, &user.Id, &user.UpdatedAt, &user.Bio)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetAllUser(req *pb.Filter) (*pb.UsersInfo, error) {
	var params []interface{}

	query := `
		select
				id, 
				user_name, 
				user_type,
				full_name
		from 
			users 
		where  
			deleted_at is null `

	if req.Limit > 0 {
		params = append(params, req.Limit)
		query += fmt.Sprintf(" and limit = $%d", len(params))
	}

	if req.Page > 0 {
		params = append(params, req.Page)
		query += fmt.Sprintf(" and offset = $%d", len(params))
	}

	rows, err := u.Db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	var users pb.UsersInfo
	for rows.Next() {
		var user pb.GetUsers1
		err := rows.Scan(&user.Id, &user.UserName, &user.UserType, &user.FullName)
		if err != nil {
			return nil, err
		}
		users.Users = append(users.Users, &user)
	}
	users.Limitpage.Limit = req.Limit
	users.Limitpage.Page = req.Page
	users.Total = int32(len(users.Users))
	return &users, nil

}

func (u *UserRepo) ValidateUserId(rep *pb.Id) (*pb.Exists, error) {
	query := `select 
	            case 
				    when id = $1  then true 
				else 
				    false 
				end 
			from 
			    users 
			where 
			    id = $1 and deletad_at is null`
	res := pb.Exists{}
	err := u.Db.QueryRow(query, rep.Id).Scan(&res.Exist)
	return &res, err
}



