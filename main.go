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

	server := api.NewServer(context.Background(), os.Getenv("ENVIRONMENT"))

	serverAddress := os.Getenv("SERVER_ADDRESS")
	fmt.Println(serverAddress)
	
	err := server.InitSuperUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("server starting...")
	
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
