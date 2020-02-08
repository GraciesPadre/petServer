package main

import (
	"log"
	"net/http"
	"os"
	"petServer/dataStore"
	"petServer/webServer"
)

func main() {
	store, err := dataStore.NewDataStore("/Users/doomer/tmp/pets.json")

	if err != nil {
		log.Printf("error : %+v", err)
		os.Exit(-1)
	}

	server, err := webServer.NewPetServer(":8080", store)

	if err != nil {
		log.Printf("Error: %+v", err)
		os.Exit(-1)
	}

	if err := server.Start(); err != http.ErrServerClosed {
		log.Printf("Error: %+v\n", err)
		os.Exit(-1)
	}
}
