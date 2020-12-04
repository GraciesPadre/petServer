package dataStore

import (
	"sync"
)

func NewDataStore(filePath string) (DataStore, error) {
	serverSettings, err := NewServerSettings(filePath)

	if err != nil {
		return nil, err
	}
	return &dataStore{serverSettings: serverSettings, petsCollection: NewPetsCollection()}, nil
}

type Loader interface {
	Load() error
}

type Storeer interface {
	Store() error
}

type DataStore interface {
	Loader
	Storeer
	AddPet(name string, breed string, age int) PetsCollection
	RemovePet(name string) PetsCollection
	AllPets() PetsCollection
	OnePet(name string) PetsCollection
}

type dataStore struct {
	serverSettings ServerSettings
	petsCollection PetsCollection
	lock           sync.RWMutex
}

func (store *dataStore) Load() error {
	store.lock.Lock()
	defer store.lock.Unlock()

	petsCollection, err := store.serverSettings.Deserialize()

	if err != nil {
		return err
	}

	store.petsCollection = petsCollection

	return nil
}

func (store *dataStore) Store() error {
	store.lock.Lock()
	defer store.lock.Unlock()

	if err := store.serverSettings.Serialize(store.petsCollection); err != nil {
		return err
	}

	return nil
}

func (store *dataStore) AddPet(name string, breed string, age int) PetsCollection {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.petsCollection.Collection[name] = Pet{Age: age, Breed: breed}

	return store.petsCollection
}

func (store *dataStore) RemovePet(name string) PetsCollection {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.petsCollection.Collection, name)

	return store.petsCollection
}

func (store *dataStore) AllPets() PetsCollection {
	store.lock.RLock()
	defer store.lock.RUnlock()

	return store.petsCollection
}

func (store *dataStore) OnePet(name string) PetsCollection {
	store.lock.RLock()
	defer store.lock.RUnlock()

	result := NewPetsCollection()

	pet := store.petsCollection.Collection[name]

	if pet.Age != 0 || len(pet.Breed) > 0 {
		result.Collection[name] = pet
	}

	return result
}
