package webServer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"petServer/dataStore"
	"strings"
	"testing"
)

func TestWritingThenReadingPetSettings(t *testing.T) {
	const settingsFilePath = "ick.json"
	const firstPetsPath = settingsFilePath
	const secondPetsPath = "poo.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	const breed1 = "breed 1"
	const breed2 = "breed 2"

	petsCollection := dataStore.NewPetsCollection()
	petsCollection.Collection[firstPetsPath] = dataStore.Pet{Age: 1, Breed: breed1}
	petsCollection.Collection[secondPetsPath] = dataStore.Pet{Age: 2, Breed: breed2}

	settings, err := dataStore.NewServerSettings(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	err = settings.Serialize(petsCollection)

	if err != nil {
		t.Fatal(err)
	}

	deserializedPets, err := settings.Deserialize()

	if err != nil {
		t.Fatal(err)
	}

	if deserializedPets.Collection[firstPetsPath].Age != 1 || deserializedPets.Collection[firstPetsPath].Breed != breed1 {
		t.Fatalf("deserializedPets.Collection[firstPetsPath].Age != 1 || deserializedPets.Collection[firstPetsPath].Breed != breed1")
	}

	if deserializedPets.Collection[secondPetsPath].Age != 2 || deserializedPets.Collection[secondPetsPath].Breed != breed2 {
		t.Fatalf("deserializedPets.Collection[secondPetsPath].Age != 2 || deserializedPets.Collection[secondPetsPath].Breed != breed2")
	}
}

func remove(fileName string) {
	_ = os.Remove(fileName)
}

func TestGettingUndefinedPet(t *testing.T) {
	settingsFilePath := "TestGettingUndefinedPet.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	store, err := dataStore.NewDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	allPets := store.AllPets()

	if len(allPets.Collection) != 0 {
		t.Fatal("got a pet when none was expected")
	}
}

func TestGettingDefinedPet(t *testing.T) {
	settingsFilePath := "TestGettingDefinedPet.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	store, err := dataStore.NewDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	const shasta = "Shasta"
	const shastaBreed = "Eskie"

	store.AddPet(shasta, shastaBreed, 9)

	petsCollection := store.OnePet(shasta)

	if len(petsCollection.Collection) != 1 {
		t.Fatalf("extected 1 pet but got %d", len(petsCollection.Collection))
	}

	pet := petsCollection.Collection[shasta]

	if pet.Breed != shastaBreed {
		t.Fatal("got the wrong breed")
	}

	if pet.Age != 9 {
		t.Fatal("got the wrong age")
	}
}

func TestAddingAndRemovingPet(t *testing.T) {
	settingsFilePath := "TestGettingDefinedPet.json"

	_ = os.Remove(settingsFilePath)
	defer remove(settingsFilePath)

	store, err := dataStore.NewDataStore(settingsFilePath)

	if err != nil {
		t.Fatal(err)
	}

	const shasta = "Shasta"
	const shastaBreed = "Eskie"

	store.AddPet(shasta, shastaBreed, 9)

	petsCollection := store.OnePet(shasta)

	if len(petsCollection.Collection) != 1 {
		t.Fatalf("extected 1 pet but got %d", len(petsCollection.Collection))
	}

	pet := petsCollection.Collection[shasta]

	if pet.Breed != shastaBreed {
		t.Fatal("got the wrong breed")
	}

	if pet.Age != 9 {
		t.Fatal("got the wrong age")
	}

	petsCollection = store.RemovePet(shasta)

	if len(petsCollection.Collection) != 0 {
		t.Fatal("pets collection shoe have no entries")
	}
}

func TestGettingEmptyCollection(t *testing.T) {
	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(mockGetHandler)

	request, err := http.NewRequest("GET", "/pet", nil)
	if err != nil {
		t.Fatal(err)
	}

	httpHandler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"pets_collection":{}}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func mockGetHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	newStore, err := dataStore.NewDataStore("./mockGetHandler.json")

	if err != nil {
		panic(err)
	}

	handler := &getHandler{newStore}

	if err := handler.HandleRequest(responseWriter, httpRequest); err != nil {
		panic(err)
	}
}

func TestGettingUndefinedPetFromEmptyCollection(t *testing.T) {
	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(mockGetHandler)

	request, err := http.NewRequest("GET", "/pet?name=Gracie", nil)
	if err != nil {
		t.Fatal(err)
	}

	httpHandler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"pets_collection":{}}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestAdding1Pet(t *testing.T) {
	const filePath = "TestAdding1Pet.json"

	defer remove(filePath)

	recorder := httptest.NewRecorder()

	store, err := dataStore.NewDataStore(filePath)

	if err != nil {
		t.Fatalf("making data store fail with error: %+v", err)
	}

	thePutHandler := &mockPutHandler{store: store}

	httpHandler := http.HandlerFunc(thePutHandler.mockPutHandlerWithDataStore)

	expected := `{"pets_collection":{"Shasta":{"age":9,"breed":"Spitz"}}}`

	request, err := http.NewRequest("PUT", "/pet", strings.NewReader(expected))
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

	theGetHandler := &mockGetHandler2{store: store}

	recorder2 := httptest.NewRecorder()
	httpHandler2 := http.HandlerFunc(theGetHandler.mockGetHandlerWithDataStore)

	request2, err := http.NewRequest("GET", "/pet?name=Shasta", nil)

	if err != nil {
		t.Fatal(err)
	}

	httpHandler2.ServeHTTP(recorder2, request2)

	if status := recorder2.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if recorder2.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder2.Body.String(), expected)
	}
}

type mockPutHandler struct {
	store dataStore.DataStore
}

func (mockHandler *mockPutHandler) mockPutHandlerWithDataStore(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	handler := &putHandler{mockHandler.store}

	if err := handler.HandleRequest(responseWriter, httpRequest); err != nil {
		panic(err)
	}
}

type mockGetHandler2 struct {
	store dataStore.DataStore
}

func (mockHandler *mockGetHandler2) mockGetHandlerWithDataStore(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	handler := &getHandler{mockHandler.store}

	if err := handler.HandleRequest(responseWriter, httpRequest); err != nil {
		panic(err)
	}
}

func TestGettingUndefinedPetHttp(t *testing.T) {
	const filePath = "TestGettingUndefinedPetHttp.json"

	defer remove(filePath)

	recorder := httptest.NewRecorder()

	store, err := dataStore.NewDataStore(filePath)

	if err != nil {
		t.Fatalf("making data store fail with error: %+v", err)
	}

	thePutHandler := &mockPutHandler{store: store}

	httpHandler := http.HandlerFunc(thePutHandler.mockPutHandlerWithDataStore)

	const petDefinition = `{"pets_collection":{"Shasta":{"age":9,"breed":"Spitz"}}}`

	request, err := http.NewRequest("PUT", "/pet", strings.NewReader(petDefinition))
	if err != nil {
		t.Fatal(err)
	}

	httpHandler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if recorder.Body.String() != petDefinition {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), petDefinition)
	}

	theGetHandler := &mockGetHandler2{store: store}

	recorder2 := httptest.NewRecorder()
	httpHandler2 := http.HandlerFunc(theGetHandler.mockGetHandlerWithDataStore)

	request2, err := http.NewRequest("GET", "/pet?name=Gracie", nil)

	if err != nil {
		t.Fatal(err)
	}

	httpHandler2.ServeHTTP(recorder2, request2)

	if status := recorder2.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	const expected = `{"pets_collection":{}}`
	if recorder2.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder2.Body.String(), expected)
	}
}

func TestDeletingPet(t *testing.T) {
	const filePath = "TestAdding1Pet.json"

	defer remove(filePath)

	recorder := httptest.NewRecorder()

	store, err := dataStore.NewDataStore(filePath)

	if err != nil {
		t.Fatalf("making data store fail with error: %+v", err)
	}

	thePutHandler := &mockPutHandler{store: store}

	httpHandler := http.HandlerFunc(thePutHandler.mockPutHandlerWithDataStore)

	petDefinition := `{"pets_collection":{"Shasta":{"age":9,"breed":"Spitz"}}}`

	request, err := http.NewRequest("PUT", "/pet", strings.NewReader(petDefinition))
	if err != nil {
		t.Fatal(err)
	}

	httpHandler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if recorder.Body.String() != petDefinition {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), petDefinition)
	}

	theGetHandler := &mockGetHandler2{store: store}

	recorder2 := httptest.NewRecorder()
	httpHandler2 := http.HandlerFunc(theGetHandler.mockGetHandlerWithDataStore)

	request2, err := http.NewRequest("GET", "/pet?name=Shasta", nil)

	if err != nil {
		t.Fatal(err)
	}

	httpHandler2.ServeHTTP(recorder2, request2)

	if status := recorder2.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if recorder2.Body.String() != petDefinition {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder2.Body.String(), petDefinition)
	}

	theDeleteHandler := &mockDeleteHandler{store: store}

	recorder3 := httptest.NewRecorder()
	httpHandler3 := http.HandlerFunc(theDeleteHandler.mockDeleteHandlerWithDataStore)

	request3, err := http.NewRequest("DELETE", "/pet?name=Shasta", strings.NewReader(petDefinition))

	if err != nil {
		t.Fatal(err)
	}

	httpHandler3.ServeHTTP(recorder3, request3)

	if status := recorder3.Code; status != http.StatusOK {
		t.Errorf("httpHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	const expected = `{"pets_collection":{}}`
	if recorder3.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder3.Body.String(), expected)
	}
}

type mockDeleteHandler struct {
	store dataStore.DataStore
}

func (mockHandler *mockDeleteHandler) mockDeleteHandlerWithDataStore(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	handler := &deleteHandler{mockHandler.store}

	if err := handler.HandleRequest(responseWriter, httpRequest); err != nil {
		panic(err)
	}
}
