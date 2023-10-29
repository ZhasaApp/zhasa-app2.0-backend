package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"zhasa2.0/api"
)

func main() {

	server := api.NewServer(context.Background())

	serverAddress := os.Getenv("SERVER_ADDRESS")
	fmt.Println(serverAddress)
	err := server.InitSuperUser()

	if err != nil {
		log.Fatal(err)
	}
	err = server.Start(serverAddress)

	fmt.Println("server started")

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
