package server

import (
	"auth_service/config"
	pb "auth_service/genproto/auth"
	"auth_service/service"
	"auth_service/storage/postgres"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServr() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg := config.Load()

	listen, err := net.Listen("tcp", "localhost:"+cfg.AUTH_SERVICE_PORT)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer listen.Close()

	authService := service.NewHadler(db)

	s := grpc.NewServer()

	pb.RegisterAuthServer(s, authService)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
