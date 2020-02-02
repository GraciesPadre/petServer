package webServer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func NewCiWebServer(port string, dataStore CiDataStore) (CiWebServer, error) {
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

	return &ciWebServer{
		port:       port,
		httpServer: nil,
		dispatcher: dispatcher,
		dataStore:  dataStore,
	}, nil
}

type CiWebServer interface {
	Start() error
	Stop(responseWriter http.ResponseWriter, httpRequest *http.Request)
	HandleTestInfo(responseWriter http.ResponseWriter, httpRequest *http.Request)
}

type ciWebServer struct {
	port       string
	httpServer *http.Server
	dispatcher Dispatcher
	dataStore  CiDataStore
	lock       sync.Mutex
}

func (server *ciWebServer) Start() error {
	_ = server.dataStore.Load()

	server.newServer()

	return server.httpServer.ListenAndServe()
}

func (server *ciWebServer) newServer() {
	server.lock.Lock()
	defer server.lock.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc("/close", server.Stop)
	mux.HandleFunc("/integrationTest", server.HandleTestInfo)

	server.httpServer = &http.Server{Addr: server.port, Handler: mux}
}

/*
curl --header "Content-Type: application/json" --request PUT --data '{"settings_collection":{"ick":{"enabled":false,"gates_ci_build":false}}}' http://localhost:8080/integrationTest
curl http://localhost:8080/integrationTest?testPath=Gracie
curl http://localhost:8080/integrationTest
curl -X DELETE http://localhost:8080/integrationTest?testPath=Shasta
curl -X PUT http://localhost:8080/close
*/
func (server *ciWebServer) HandleTestInfo(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if err := server.dispatcher.HandleRequest(responseWriter, httpRequest); err != nil {
		responseWriter.WriteHeader(500)
	}
}

func (server *ciWebServer) Stop(responseWriter http.ResponseWriter, httpRequest *http.Request) {
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
