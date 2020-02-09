package dataStore

import (
	"testing"
)

func TestLoadingDataStore(t *testing.T) {
	const fileName = "TestLoadingDataStore.json"

	defer nukeFile(fileName)

	_ = store3Pets(t, fileName)

	store, err := NewDataStore(fileName)

	if err != nil {
		t.Error(err)
	}

	if err := store.Load(); err != nil {
		t.Error(err)
	}
}

func TestGettingNonExistentPet(t *testing.T) {
	const fileName = "TestGettingNonExistentPet.json"

	defer nukeFile(fileName)

	_ = store3Pets(t, fileName)

	store, err := NewDataStore(fileName)

	if err != nil {
		t.Error(err)
	}

	if err := store.Load(); err != nil {
		t.Error(err)
	}

	petCollection := store.OnePet("noName")

	if len(petCollection.Collection) > 0 {
		t.Error("got pet where none was expected")
	}
}

func TestGetting1Pet(t *testing.T) {
	const fileName = "TestGetting1Pet.json"

	defer nukeFile(fileName)

	_ = store3Pets(t, fileName)

	store, err := NewDataStore(fileName)

	if err != nil {
		t.Error(err)
	}

	if err := store.Load(); err != nil {
		t.Error(err)
	}

	petCollection := store.OnePet(buttons)

	if len(petCollection.Collection) != 1 {
		t.Errorf("expected to get 1 pet, but got %d", len(petCollection.Collection))
	}

	pet := petCollection.Collection[buttons]

	if pet.Breed != buttonsBreed || pet.Age != buttonsAge {
		t.Error("got the wrong pet")
	}
}

func TestGettingAllPets(t *testing.T) {
	const fileName = "TestGettingAllPets.json"

	defer nukeFile(fileName)

	_ = store3Pets(t, fileName)

	store, err := NewDataStore(fileName)

	if err != nil {
		t.Error(err)
	}

	if err := store.Load(); err != nil {
		t.Error(err)
	}

	pets := store.AllPets()

	if len(pets.Collection) != 3 {
		t.Errorf("expected 3 pets, got %d", len(pets.Collection))
	}

	pet := pets.Collection[buttons]
	if pet.Age != buttonsAge || pet.Breed != buttonsBreed {
		t.Errorf("collection does not contain pet named %s", buttons)
	}

	pet = pets.Collection[shasta]
	if pet.Age != shastaAge || pet.Breed != shastaBreed {
		t.Errorf("collection does not contain pet named %s", shasta)
	}

	pet = pets.Collection[gracie]
	if pet.Age != gracieAge || pet.Breed != gracieBreed {
		t.Errorf("collection does not contain pet named %s", gracie)
	}
}

func TestStoringPets(t *testing.T) {
	const fileName = "TestStoringPets.json"

	defer nukeFile(fileName)

	store, err := NewDataStore(fileName)

	if err != nil {
		t.Error(err)
	}

	pets := store.AllPets()

	if len(pets.Collection) != 0 {
		t.Errorf("expected 0 pets, got %d", len(pets.Collection))
	}

	const twitch = "Twitch"
	const twitchBreed = "Dutch Belted"
	const twitchAge = 13

	store.AddPet(twitch, twitchBreed, twitchAge)

	if err := store.Store(); err != nil {
		t.Error(err)
	}

	store2, err := NewDataStore(fileName)

	if err != nil {
		t.Error(err)
	}

	if err := store2.Load(); err != nil {
		t.Error(nil)
	}

	pets2 := store2.AllPets()

	if len(pets2.Collection) != 1 {
		t.Errorf("expected 1 pet, got %d", len(pets2.Collection))
	}

	pet := pets2.Collection[twitch]

	if pet.Age != twitchAge || pet.Breed != twitchBreed {
		t.Error("got incorrect pet info")
	}
}

func TestLoadingFromNonExistentFile(t *testing.T) {
	const fileName = "TestLoadingFromNonExistentFile.json"

	_, err := NewDataStore(fileName)

	if err != nil {
		t.Fatal("trying to read non-existent file should not error")
	}
}

func TestSavingEmptyCollection(t *testing.T) {
	const fileName = "TestSavingEmptyCollection.json"

	defer nukeFile(fileName)

	newStore, err := NewDataStore(fileName)

	if err != nil {
		t.Fatal("trying to read non-existent file should not error")
	}

	err = newStore.Store()

	if err != nil {
		t.Fatal("storing empty collection should not error")
	}
}
