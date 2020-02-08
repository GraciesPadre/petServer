package webServer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"petServer/dataStore"
	"sync"
)

func NewPetServer(port string, dataStore dataStore.DataStore) (PetServer, error) {
	if len(port) == 0 {
		return nil, fmt.Errorf("port may not be empty")
	}

	if dataStore == nil {
		return nil, fmt.Errorf("dataStore may not be nil")
	}

	dispatcher, err := NewDispatcher(dataStore)

	if err != nil {
		return nil, err
	}

	return &petServer{
		port:       port,
		httpServer: nil,
		dispatcher: dispatcher,
		dataStore:  dataStore,
	}, nil
}

type PetServer interface {
	Start() error
	Stop(responseWriter http.ResponseWriter, httpRequest *http.Request)
	HandlePetInfo(responseWriter http.ResponseWriter, httpRequest *http.Request)
}

type petServer struct {
	port       string
	httpServer *http.Server
	dispatcher Dispatcher
	dataStore  dataStore.DataStore
	lock       sync.Mutex
}

func (server *petServer) Start() error {
	_ = server.dataStore.Load()

	server.newServer()

	return server.httpServer.ListenAndServe()
}

func (server *petServer) newServer() {
	server.lock.Lock()
	defer server.lock.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc("/close", server.Stop)
	mux.HandleFunc("/pet", server.HandlePetInfo)

	server.httpServer = &http.Server{Addr: server.port, Handler: mux}
}

/*
curl --header "Content-Type: application/json" -X PUT --data '{"pets_collection":{"Buttons":{"age":2,"breed":"Terrier"},"Gracie":{"age":9,"breed":"Spitz"},"Shasta":{"age":9,"breed":"Spitz"}}}' http://localhost:8080/pet
curl http://localhost:8080/pet
curl http://localhost:8080/pet?name=Buttons
curl -X DELETE http://localhost:8080/pet?name=Shastas
curl -X PUT http://localhost:8080/close
*/
func (server *petServer) HandlePetInfo(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if err := server.dispatcher.HandleRequest(responseWriter, httpRequest); err != nil {
		responseWriter.WriteHeader(500)
	}
}

func (server *petServer) Stop(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method == "PUT" {
		server.lock.Lock()
		defer server.lock.Unlock()

		if err := server.dataStore.Store(); err != nil {
			log.Printf("storing failed with error: %+v\n", err)
		}

		if err := server.httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %+v\n", err)
		}

		if err := server.httpServer.Close(); err != nil {
			log.Printf("HTTP server Close: %+v\n", err)
		}

		server.httpServer = nil
	}
}
