package storage

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

// ListIDs returns a list of all existing paste IDs
func (driver *FileDriver) ListIDs() ([]string, error) {
	// Define the IDs slice
	var ids []string

	// Fill the IDs slice
	err := filepath.Walk(driver.filePath, func(path string, info os.FileInfo, err error) error {
		// Check if a walking error occurred
		if err != nil {
			return err
		}

		// Decode the file name
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSuffix(info.Name(), ".json"))
		if err != nil {
			return err
		}

		// Append the ID to the IDs slice
		ids = append(ids, string(decoded))
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Return the IDs slice
	return ids, nil
}

// Get loads a paste
func (driver *FileDriver) Get(id string) (*pastes.Paste, error) {
	// Read the file
	id = base64.StdEncoding.EncodeToString([]byte(id))
	data, err := ioutil.ReadFile(filepath.Join(driver.filePath, id+".json"))
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
	id := base64.StdEncoding.EncodeToString([]byte(paste.ID))
	file, err := os.Create(filepath.Join(driver.filePath, id+".json"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data into the file
	_, err = file.Write(jsonBytes)
	return err
}

// Delete deletes a paste
func (driver *FileDriver) Delete(id string) error {
	id = base64.StdEncoding.EncodeToString([]byte(id))
	return os.Remove(filepath.Join(driver.filePath, id+".json"))
}
