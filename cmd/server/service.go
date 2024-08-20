package server

import (
	"auth_service/config"
	pb "auth_service/genproto/auth"
	"auth_service/service"
	"auth_service/storage/postgres"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServr() {
	fmt.Println("------------------------")
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg := config.Load()
	fmt.Println("-----------------------+")

	listen, err := net.Listen("tcp", "localhost:"+cfg.AUTH_SERVICE_PORT)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("------------------------+-")

	defer listen.Close()

	authService := service.NewHadler(db)

	s := grpc.NewServer()

	fmt.Println("------------------------+-+")
	pb.RegisterAuthServer(s, authService)
	fmt.Println("------------------------+-+-")
	
	if err := s.Serve(listen); err != nil {
	fmt.Println("------------------------+-+-+")
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("+++++++++++++++++++++++++++++++++++")
	fmt.Println("grpc listen :7777 ...")
}
