package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"zhasa2.0/api"
	"zhasa2.0/pkg/notify"
)

func main() {

	client, err := notify.NewFirebaseClient("config/firebase.json")
	if err != nil {
		log.Fatal("cannot initialize firebase client", err)
	}

	server := api.NewServer(context.Background(), os.Getenv("ENVIRONMENT"))
	server.SetFBClient(client)

	serverAddress := os.Getenv("SERVER_ADDRESS")
	fmt.Println(serverAddress)

	err = server.InitSuperUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("server starting...")

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
