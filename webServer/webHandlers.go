package webServer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"petServer/dataStore"
)

type HttpRequestHandler interface {
	HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error
}

type PutHandler interface {
	HttpRequestHandler
	HandlePut(responseWriter http.ResponseWriter, httpRequest *http.Request) error
}

type putHandler struct {
	dataStore dataStore.DataStore
}

func (handler *putHandler) HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	return handler.HandlePut(responseWriter, httpRequest)
}

func (handler *putHandler) HandlePut(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	contentLength := httpRequest.ContentLength
	body := make([]byte, contentLength)
	_, err := httpRequest.Body.Read(body)

	if err != nil && err != io.EOF {
		responseWriter.WriteHeader(500)
		return err
	}

	var settingsCollection dataStore.PetsCollection
	err = json.Unmarshal(body, &settingsCollection)

	if err != nil {
		responseWriter.WriteHeader(500)
		return err
	}

	for name, settings := range settingsCollection.Collection {
		handler.dataStore.AddPet(name, settings.Breed, settings.Age)
	}

	return getAllSettings(handler.dataStore, responseWriter)
}

func getAllSettings(dataStore dataStore.DataStore, responseWriter http.ResponseWriter) error {
	pets := dataStore.AllPets()

	result, err := json.Marshal(pets)

	if err != nil {
		responseWriter.WriteHeader(500)
		return err
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	_, _ = responseWriter.Write(result)

	return nil
}

type GetHandler interface {
	HttpRequestHandler
	HandleGet(responseWriter http.ResponseWriter, httpRequest *http.Request) error
}

type getHandler struct {
	dataStore dataStore.DataStore
}

func (handler *getHandler) HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	return handler.HandleGet(responseWriter, httpRequest)
}

func (handler *getHandler) HandleGet(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	name := httpRequest.URL.Query().Get("name")

	if len(name) == 0 {
		return handler.handleGetAllSettings(responseWriter)
	}

	return handler.handleGetOneSetting(responseWriter, name)
}

func (handler *getHandler) handleGetAllSettings(responseWriter http.ResponseWriter) error {
	return getAllSettings(handler.dataStore, responseWriter)
}

func (handler *getHandler) handleGetOneSetting(responseWriter http.ResponseWriter, name string) error {
	pet := handler.dataStore.OnePet(name)

	result, err := json.Marshal(pet)

	if err != nil {
		responseWriter.WriteHeader(500)
		return err
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	_, _ = responseWriter.Write(result)

	return nil
}

type DeleteHandler interface {
	HttpRequestHandler
	HandleDelete(responseWriter http.ResponseWriter, httpRequest *http.Request) error
}

type deleteHandler struct {
	dataStore dataStore.DataStore
}

func (handler *deleteHandler) HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	return handler.HandleDelete(responseWriter, httpRequest)
}

func (handler *deleteHandler) HandleDelete(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	name := httpRequest.URL.Query().Get("name")

	if len(name) == 0 {
		responseWriter.WriteHeader(400)
		return fmt.Errorf("testPath not found in parameters")
	}

	settingsCollection := handler.dataStore.RemovePet(name)

	result, err := json.Marshal(settingsCollection)

	if err != nil {
		responseWriter.WriteHeader(500)
		return err
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	_, _ = responseWriter.Write(result)

	return nil
}
