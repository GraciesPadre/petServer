package dataStore

func NewPetsCollection() PetsCollection {
	petsCollection := PetsCollection{}
	petsCollection.Collection = make(map[string]Pet)
	return petsCollection
}

type Pet struct {
	Age   int    `json:"age"`
	Breed string `json:"breed"`
}

type PetsCollection struct {
	Collection map[string]Pet `json:"pets_collection"`
}
