package dataStore

import (
	"os"
	"testing"
)

const shasta = "Shasta"
const shastaAge = 9
const shastaBreed = "Spitz"

const gracie = "Gracie"
const gracieAge = 9
const gracieBreed = "Spitz"

const buttons = "Buttons"
const buttonsAge = 2
const buttonsBreed = "Terrier"

func TestStoringOnePet(t *testing.T) {
	filePath := "TestStoringOnePet.json"
	defer nukeFile(filePath)

	petsCollection := NewPetsCollection()
	petsCollection.Collection[shasta] = Pet{Age: shastaAge, Breed: shastaBreed}

	settings, err := NewServerSettings(filePath)

	if err != nil {
		t.Error(err)
	}

	err = settings.Serialize(petsCollection)

	if err != nil {
		t.Error(err)
	}

	deserializedCollection, err := settings.Deserialize()

	if err != nil {
		t.Error(err)
	}

	if len(deserializedCollection.Collection) != 1 {
		t.Errorf("collection had %d enetries, expected 1", len(deserializedCollection.Collection))
	}

	pet := deserializedCollection.Collection[shasta]

	if pet.Age != shastaAge {
		t.Errorf("Age is %d, expected %d", pet.Age, shastaAge)
	}

	if pet.Breed != shastaBreed {
		t.Errorf("Breed is %s, expected %s", pet.Breed, shastaBreed)
	}
}

func nukeFile(filePath string) {
	_ = os.Remove(filePath)
}

func TestStoring3Pets(t *testing.T) {
	const filePath = "TestStoring3Pets.json"

	defer nukeFile(filePath)

	_ = store3Pets(t, filePath)

	settings, err := NewServerSettings(filePath)

	if err != nil {
		t.Error(err)
	}

	deserializedCollection, err := settings.Deserialize()

	if err != nil {
		t.Error(err)
	}

	if len(deserializedCollection.Collection) != 3 {
		t.Errorf("collection had %d enetries, expected 3", len(deserializedCollection.Collection))
	}

	pet := deserializedCollection.Collection[shasta]

	if pet.Age != shastaAge {
		t.Errorf("Age is %d, expected %d", pet.Age, shastaAge)
	}

	if pet.Breed != shastaBreed {
		t.Errorf("Breed is %s, expected %s", pet.Breed, shastaBreed)
	}

	pet = deserializedCollection.Collection[gracie]

	if pet.Age != gracieAge {
		t.Errorf("Age is %d, expected %d", pet.Age, gracieAge)
	}

	if pet.Breed != gracieBreed {
		t.Errorf("Breed is %s, expected %s", pet.Breed, gracieBreed)
	}

	pet = deserializedCollection.Collection[buttons]

	if pet.Age != buttonsAge {
		t.Errorf("Age is %d, expected %d", pet.Age, buttonsAge)
	}

	if pet.Breed != buttonsBreed {
		t.Errorf("Breed is %s, expected %s", pet.Breed, buttonsBreed)
	}
}

func store3Pets(t *testing.T, filePath string) PetsCollection {
	petsCollection := NewPetsCollection()
	petsCollection.Collection[shasta] = Pet{Age: shastaAge, Breed: shastaBreed}
	petsCollection.Collection[gracie] = Pet{Age: gracieAge, Breed: gracieBreed}
	petsCollection.Collection[buttons] = Pet{Age: buttonsAge, Breed: buttonsBreed}

	settings, err := NewServerSettings(filePath)

	if err != nil {
		t.Error(err)
	}

	err = settings.Serialize(petsCollection)

	if err != nil {
		t.Error(err)
	}

	return petsCollection
}
