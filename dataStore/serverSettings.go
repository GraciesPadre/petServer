package dataStore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func NewServerSettings(settingsFilePath string) (ServerSettings, error) {
	if len(settingsFilePath) == 0 {
		return nil, fmt.Errorf("settingsFilePath may not be empty")
	}

	return &serverSettings{settingsFilePath: settingsFilePath}, nil
}

type ServerSettings interface {
	Serialize(petsCollection PetsCollection) error
	Deserialize() (PetsCollection, error)
}

type serverSettings struct {
	settingsFilePath string
}

func (settings *serverSettings) Serialize(petsCollection PetsCollection) error {
	_, err := os.Stat(settings.settingsFilePath)

	var file *os.File = nil

	if os.IsNotExist(err) {
		file, err = os.Create(settings.settingsFilePath)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(settings.settingsFilePath, os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if file == nil {
		return fmt.Errorf("%s in not open for writing", settings.settingsFilePath)
	} else {
		defer finalize(file)
	}

	serializedSettings, err := json.Marshal(petsCollection)
	if err != nil {
		return err
	}

	_, err = file.Write(serializedSettings)

	if err != nil {
		return err
	}

	return nil
}

func finalize(file *os.File) {
	_ = file.Close()
}

func (settings *serverSettings) Deserialize() (PetsCollection, error) {
	fileData, err := ioutil.ReadFile(settings.settingsFilePath)

	if err != nil {
		return PetsCollection{}, err
	}

	petsCollection := PetsCollection{}
	err = json.Unmarshal(fileData, &petsCollection)

	if err != nil {
		return PetsCollection{}, err
	}

	return petsCollection, nil
}
