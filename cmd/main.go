package main

import (
	"auth_service/api"
	"auth_service/cmd/server"
	"auth_service/config"
	"auth_service/storage/postgres"
	"log"
	"sync"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg := config.Load()
	
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		server.RunServr()
	}()
	router := api.NewRouter(db)
	router.Run(cfg.HTTP_PORT)

	wg.Wait()

}
