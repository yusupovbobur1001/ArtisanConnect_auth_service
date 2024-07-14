package main

import (
	"auth_service/config"
	"auth_service/service"
	"auth_service/storage/postgres"
	"log"
	"net"	
	pb "auth_service/genproto/auth"
	"google.golang.org/grpc"
)



func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Load()

	listen, err := net.Listen("tcp", cfg.AUTH_SERVICE_PORT)
	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()

	authService := service.NewHadler(db)

	s := grpc.NewServer()

	pb.RegisterAuthServer(s, authService)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}