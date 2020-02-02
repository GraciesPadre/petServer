package webServer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"petServer/dataStore"
	"strings"
	"testing"
)

func TestWritingThenReadingIntegrationTestSettings(t *testing.T) {
	const settingsFilePath = "ick.json"
	const firstIntegrationTestPath = settingsFilePath
	const secondIntegrationTestPath = "poo.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	integrationTestSettingsCollection := NewIntegrationTestSettingsCollection()
	integrationTestSettingsCollection.SettingsCollection[firstIntegrationTestPath] = IntegrationTestSettings{Enabled: true, GatesCIBuild: true}
	integrationTestSettingsCollection.SettingsCollection[secondIntegrationTestPath] = IntegrationTestSettings{Enabled: false, GatesCIBuild: false}

	settings, err := dataStore.NewServerSettings(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	err = settings.Serialize(integrationTestSettingsCollection)

	if err != nil {
		t.Fatal(err)
	}

	deserializedIntegrationTestSettings, err := settings.Deserialize()

	if err != nil {
		t.Fatal(err)
	}

	if deserializedIntegrationTestSettings.SettingsCollection[firstIntegrationTestPath].Enabled == false || deserializedIntegrationTestSettings.SettingsCollection[firstIntegrationTestPath].GatesCIBuild == false {
		t.Fatalf("deserializedIntegrationTestSettings.SettingsCollection[firstIntegrationTestPath].Enabled == false || deserializedIntegrationTestSettings.SettingsCollection[firstIntegrationTestPath].GatesCIBuild == false")
	}

	if deserializedIntegrationTestSettings.SettingsCollection[secondIntegrationTestPath].Enabled == true || deserializedIntegrationTestSettings.SettingsCollection[secondIntegrationTestPath].GatesCIBuild == true {
		t.Fatalf("deserializedIntegrationTestSettings.SettingsCollection[secondIntegrationTestPath].Enabled == true || deserializedIntegrationTestSettings.SettingsCollection[secondIntegrationTestPath].GatesCIBuild")
	}
}

func remove(fileName string) {
	_ = os.Remove(fileName)
}

func TestGettingUndefinedTestState(t *testing.T) {
	settingsFilePath := "TestGettingUndefinedTestState.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	server, err := NewCiDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	enabled := server.IsEnabled(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}

	enabled = server.IsGating(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}
}

func TestGettingDefinedTestState(t *testing.T) {
	settingsFilePath := "TestGettingDefinedTestState.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	server, err := NewCiDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	server.EnableIntegrationTest(settingsFilePath)
	server.EnableGating(settingsFilePath)

	enabled := server.IsEnabled(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}

	enabled = server.IsGating(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}
}

func TestResettingDefinedTestState(t *testing.T) {
	settingsFilePath := "TestResettingDefinedTestState.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	server, err := NewCiDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	server.EnableIntegrationTest(settingsFilePath)
	server.EnableGating(settingsFilePath)

	enabled := server.IsEnabled(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}

	enabled = server.IsGating(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}

	server.DisableIntegrationTest(settingsFilePath)
	server.DisableGating(settingsFilePath)

	enabled = server.IsEnabled(settingsFilePath)

	if enabled {
		t.Fatalf("enabled was true")
	}

	enabled = server.IsGating(settingsFilePath)

	if enabled {
		t.Fatalf("enabled was true")
	}
}

func TestLoadingFromNonExistentFile(t *testing.T) {
	settingsFilePath := "TestLoadingFromNonExistentFile.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	server, err := NewCiDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	enabled := server.IsEnabled(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}

	enabled = server.IsGating(settingsFilePath)

	if !enabled {
		t.Fatalf("enabled was false")
	}
}

func TestSavingEmptySettings(t *testing.T) {
	settingsFilePath := "TestSavingEmptySettings.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	server, err := NewCiDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	err = server.Store()

	if err != nil {
		t.Fatal(err)
	}

	_, err = NewCiDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGettingEmptySettings(t *testing.T) {
	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(mockGetHandler)

	request, err := http.NewRequest("GET", "/integrationTest", nil)
	if err != nil {
		t.Fatal(err)
	}

	httpHandler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"settings_collection":{}}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func mockGetHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	dataStore, err := NewCiDataStore("./mockGetHandler.json")

	if err != nil {
		panic(err)
	}

	handler := &getHandler{dataStore}

	if err := handler.HandleRequest(responseWriter, httpRequest); err != nil {
		panic(err)
	}
}

func TestGettingUndefinedSettingFromEmptyCollection(t *testing.T) {
	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(mockGetHandler)

	request, err := http.NewRequest("GET", "/integrationTest?testPath=Gracie", nil)
	if err != nil {
		t.Fatal(err)
	}

	httpHandler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"settings_collection":{"Gracie":{"enabled":true,"gates_ci_build":true}}}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestAddingTestSetting(t *testing.T) {
	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(mockPutHandler)

	expected := `{"settings_collection":{"ick":{"enabled":false,"gates_ci_build":false}}}`

	request, err := http.NewRequest("PUT", "/integrationTest", strings.NewReader(expected))
	if err != nil {
		t.Fatal(err)
	}

	httpHandler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func mockPutHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	dataStore, err := NewCiDataStore("./mockPutHandler.json")

	if err != nil {
		panic(err)
	}

	handler := &putHandler{dataStore}

	if err := handler.HandleRequest(responseWriter, httpRequest); err != nil {
		panic(err)
	}
}
