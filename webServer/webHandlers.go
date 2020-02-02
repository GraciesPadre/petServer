package webServer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpRequestHandler interface {
	HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error
}

type PutHandler interface {
	HttpRequestHandler
	HandlePut(responseWriter http.ResponseWriter, httpRequest *http.Request) error
}

type putHandler struct {
	ciDataStore CiDataStore
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

	var settingsCollection IntegrationTestSettingsCollection
	err = json.Unmarshal(body, &settingsCollection)

	if err != nil {
		responseWriter.WriteHeader(500)
		return err
	}

	for testPath, settings := range settingsCollection.SettingsCollection {
		if settings.GatesCIBuild {
			handler.ciDataStore.EnableGating(testPath)
		} else {
			handler.ciDataStore.DisableGating(testPath)
		}

		if settings.Enabled {
			handler.ciDataStore.EnableIntegrationTest(testPath)
		} else {
			handler.ciDataStore.DisableIntegrationTest(testPath)
		}
	}

	return getAllSettings(handler.ciDataStore, responseWriter)
}

func getAllSettings(dataStore CiDataStore, responseWriter http.ResponseWriter) error {
	settings := dataStore.ReportAllTests()

	result, err := json.Marshal(settings)

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
	ciDataStore CiDataStore
}

func (handler *getHandler) HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	return handler.HandleGet(responseWriter, httpRequest)
}

func (handler *getHandler) HandleGet(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	testPath := httpRequest.URL.Query().Get("testPath")

	if len(testPath) == 0 {
		return handler.handleGetAllSettings(responseWriter)
	}

	return handler.handleGetOneSetting(responseWriter, testPath)
}

func (handler *getHandler) handleGetAllSettings(responseWriter http.ResponseWriter) error {
	return getAllSettings(handler.ciDataStore, responseWriter)
}

func (handler *getHandler) handleGetOneSetting(responseWriter http.ResponseWriter, testPath string) error {
	isEnabled := handler.ciDataStore.IsEnabled(testPath)
	isGating := handler.ciDataStore.IsGating(testPath)

	settingsMap := map[string]IntegrationTestSettings{
		testPath: {
			Enabled:      isEnabled,
			GatesCIBuild: isGating,
		},
	}

	settingsCollection := IntegrationTestSettingsCollection{SettingsCollection: settingsMap}

	result, err := json.Marshal(settingsCollection)

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
	ciDataStore CiDataStore
}

func (handler *deleteHandler) HandleRequest(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	return handler.HandleDelete(responseWriter, httpRequest)
}

func (handler *deleteHandler) HandleDelete(responseWriter http.ResponseWriter, httpRequest *http.Request) error {
	testPath := httpRequest.URL.Query().Get("testPath")

	if len(testPath) == 0 {
		responseWriter.WriteHeader(400)
		return fmt.Errorf("testPath not found in parameters")
	}

	settingsCollection := handler.ciDataStore.RemoveTestSetting(testPath)

	result, err := json.Marshal(settingsCollection)

	if err != nil {
		responseWriter.WriteHeader(500)
		return err
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	_, _ = responseWriter.Write(result)

	return nil
}
