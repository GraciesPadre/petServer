package main

import (
	"ciGatingServer/webServer"
	"log"
	"net/http"
	"os"
)

func main() {
	dataStore, err := ciGatingServer.NewCiDataStore("/dataStore/ciGatingServerSettings.json")

	if err != nil {
		log.Printf("error : %+v", err)
		os.Exit(-1)
	}

	webServer, err := ciGatingServer.NewCiWebServer(":8080", dataStore)

	if err != nil {
		log.Printf("Error: %+v", err)
		os.Exit(-1)
	}

	if err := webServer.Start(); err != http.ErrServerClosed {
		log.Printf("Error: %+v\n", err)
		os.Exit(-1)
	}
}
