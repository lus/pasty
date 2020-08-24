package storage

import (
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/bwmarrin/snowflake"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileDriver represents the file storage driver
type FileDriver struct {
	filePath string
}

// Initialize initializes the file storage driver
func (driver *FileDriver) Initialize() error {
	driver.filePath = env.Get("STORAGE_FILE_PATH", "./data")
	return os.MkdirAll(driver.filePath, os.ModePerm)
}

// Terminate terminates the file storage driver (does nothing, because the file storage driver does not need any termination)
func (driver *FileDriver) Terminate() error {
	return nil
}

// Get loads a paste
func (driver *FileDriver) Get(id snowflake.ID) (*pastes.Paste, error) {
	// Read the file
	data, err := ioutil.ReadFile(filepath.Join(driver.filePath, id.String()+".json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal the file into a paste
	paste := new(pastes.Paste)
	err = json.Unmarshal(data, &paste)
	if err != nil {
		return nil, err
	}
	return paste, nil
}

// Save saves a paste
func (driver *FileDriver) Save(paste *pastes.Paste) error {
	// Marshal the paste
	jsonBytes, err := json.Marshal(paste)
	if err != nil {
		return err
	}

	// Create the file to save the paste to
	file, err := os.Create(filepath.Join(driver.filePath, paste.ID.String()+".json"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data into the file
	_, err = file.Write(jsonBytes)
	return err
}

// Delete deletes a paste
func (driver *FileDriver) Delete(id snowflake.ID) error {
	return os.Remove(filepath.Join(driver.filePath, id.String()+".json"))
}
