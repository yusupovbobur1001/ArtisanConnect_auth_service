package postgres

import (
	pb "auth_service/genproto/auth"
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	auth := NewUserRepo(db)

	req := pb.User{
		UserName: "47d9a26d-7053-4c98-9171-5cb7c7a132d1",
		Email:    "email@gmail.com",
		Password: "pass",
		FullName: "Yusupov Bobur",
		UserType: "vendor",
		Bio:      "duradgor",
	}

	err = auth.Register(&req)
	if err != nil {
		fmt.Println("error: ", err, "++++++++++++++++++++++++++++++++++")
		panic(err)
	}
}

func TestUpdate(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	auth := NewUserRepo(db)

	id := "7808c0df-7c25-433e-83bb-af99affa62c6"

	user := pb.UserUpdate{
		UserName: "bobur",
		FullName: "bobur yusupov",
		UserType: "vendor",
		Bio:      "trfgtsfdsdssdfsadfsadfasdfadsfsadyuhju",
		Id: id,
	}

	resp, err := auth.UpdateUser(&user)
	if err != nil {
		fmt.Println(err, "+++++++++++++++++++++++++++++++++++++")
		panic(err)
	}

	fmt.Println(resp)
}

func TestDelete(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	id := "cbac7e32-3de3-4379-9a6d-36b1d1f3a506"

	auth := NewUserRepo(db)

	mes, err := auth.DeleteUser(&pb.Id{Id: id})
	fmt.Println(mes)
	if err != nil {
		panic(err)
	}
}

func TestGetIdUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	u := NewUserRepo(db)
	id := "4d1a0425-054b-4b38-a9a2-633f58e36650"

	p, err := u.GetByIdUser(&pb.Id{Id: id})
	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}

func TestGetAllUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	u := NewUserRepo(db)

	filter := pb.Filter{
		Page:  1,
		Limit: 10,
	}

	p, err := u.GetAllUser(&filter)
	if err != nil {
		panic(err)
	}

	for _, k := range p.Users {
		if len(k.FullName) == 0 || len(k.UserName) == 0 || len(k.UserType) == 0 {

			fmt.Println("malumot, tfg")
		}
	}

	fmt.Println(p)
}
