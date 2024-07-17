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
		) values($1, $2, $3, $4, $5, $6, $7) `

	_, err := u.Db.Exec(query, req.UserName, req.Email, req.Password,
		req.FullName, req.UserType, req.Bio, time.Now())

	if err != nil {
		return fmt.Errorf("register metod: %v", err)
	}

	return nil
}

func (u *UserRepo) Login(req *pb.UserLogin) (*model.User, error) {
	fmt.Println("+++++++++++++++++++++++++++")
	user := model.User{}
	fmt.Println("Email:", req.Email, "Password:", req.Password)
	query := `
		select 
			user_name, full_name, user_type, bio
		from 
			users
		where email=$1 and password=$2`

	err := u.Db.QueryRow(query, req.Email, req.Password).Scan(
				&user.UserName, 
				&user.FullName, 
				&user.UserType,  
				&user.Bio,
		)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("-------------------------------")
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
		UPDATE 
			users
		SET
			user_name = $1, 
			full_name = $2, 
			bio = $3, 
			user_type = $4, 
			updated_at = $6
		WHERE 
			id = $5
	`

	_, err := u.Db.Exec(query, req.UserName, req.FullName, req.Bio, req.UserType, req.Id, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	resp, err := u.GetByIdUser(&pb.Id{Id: req.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
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
		SELECT 
			user_name, email, password, full_name, id, updated_at, bio
		FROM 
			users
		WHERE 
			id = $1 AND deleted_at IS NULL
	`

	err := u.Db.QueryRow(query, req.Id).Scan(
		&user.UserName, 
		&user.Email, 
		&user.Password,
		&user.FullName, 
		&user.Id, 
		&user.UpdatedAt, 
		&user.Bio,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %s", req.Id)
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	return &user, nil
}


func (u *UserRepo) GetAllUser(req *pb.Filter) (*pb.UsersInfo, error) {
	var params []interface{}
	query := `
		SELECT
			id, 
			user_name, 
			user_type,
			full_name
		FROM 
			users 
		WHERE  
			deleted_at IS NULL`

	if req.Limit > 0 {
		params = append(params, req.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(params))
	}

	if req.Page > 0 {
		offset := (req.Page - 1) * req.Limit
		params = append(params, offset)
		query += fmt.Sprintf(" OFFSET $%d", len(params))
	}

	rows, err := u.Db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var users pb.UsersInfo
	for rows.Next() {
		var user pb.GetUsers1
		err := rows.Scan(&user.Id, &user.UserName, &user.UserType, &user.FullName)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users.Users = append(users.Users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
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
