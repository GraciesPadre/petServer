package webServer

import (
	"fmt"
	"net/http"
	"petServer/dataStore"
)

func NewDispatcher(dataStore dataStore.DataStore) (Dispatcher, error) {
	if dataStore == nil {
		return nil, fmt.Errorf("dataStore may not be nil")
	}
	putHandlers := []HttpRequestHandler{&putHandler{dataStore: dataStore}}
	getHandlers := []HttpRequestHandler{&getHandler{dataStore: dataStore}}
	deleteHandlers := []HttpRequestHandler{&deleteHandler{dataStore: dataStore}}

	return &dispatcher{putHandlers: putHandlers, getHandlers: getHandlers, deleteHandlers: deleteHandlers}, nil
}

type Dispatcher interface {
	AddPutHandler(postHandler PutHandler)
	AddGetHandler(getHandler GetHandler)
	AddDeleteHandler(deleteHandler DeleteHandler)
	HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error
}

type dispatcher struct {
	putHandlers    []HttpRequestHandler
	getHandlers    []HttpRequestHandler
	deleteHandlers []HttpRequestHandler
}

func (dispatcher *dispatcher) AddPutHandler(putHandler PutHandler) {
	dispatcher.putHandlers = append(dispatcher.putHandlers, putHandler)
}

func (dispatcher *dispatcher) AddGetHandler(getHandler GetHandler) {
	dispatcher.getHandlers = append(dispatcher.getHandlers, getHandler)
}

func (dispatcher *dispatcher) AddDeleteHandler(deleteHandler DeleteHandler) {
	dispatcher.deleteHandlers = append(dispatcher.deleteHandlers, deleteHandler)
}

func (dispatcher *dispatcher) HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	switch httpRequest.Method {
	case "PUT":
		return dispatcher.handleRequest(responseWriter, httpRequest, dispatcher.putHandlers)
	case "GET":
		return dispatcher.handleRequest(responseWriter, httpRequest, dispatcher.getHandlers)
	case "DELETE":
		return dispatcher.handleRequest(responseWriter, httpRequest, dispatcher.deleteHandlers)
	default:
		return fmt.Errorf("do not know how to handle request of type: %s", httpRequest.Method)
	}
}

func (dispatcher *dispatcher) handleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request, httpRequestHandlers []HttpRequestHandler) error {
	for _, handler := range httpRequestHandlers {
		if err := handler.HandleRequest(responseWriter, httpRequest); err != nil {
			return err
		}
	}

	return nil
}
