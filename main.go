package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
	"os"
	"zhasa2.0/api"
)

func main() {

	server := api.NewServer(context.Background())

	serverAddress := os.Getenv("SERVER_ADDRESS")
	err := server.InitSuperUser()
	if err != nil {
		log.Fatal("cannot init super user")
	}
	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
