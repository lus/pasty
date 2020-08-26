package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io/ioutil"
	"strings"
)

// S3Driver represents the AWS S3 storage driver
type S3Driver struct {
	client *minio.Client
	bucket string
}

// Initialize initializes the AWS S3 storage driver
func (driver *S3Driver) Initialize() error {
	client, err := minio.New(env.Get("STORAGE_S3_ENDPOINT", ""), &minio.Options{
		Creds:  credentials.NewStaticV4(env.Get("STORAGE_S3_ACCESS_KEY_ID", ""), env.Get("STORAGE_S3_SECRET_ACCESS_KEY", ""), env.Get("STORAGE_S3_SECRET_TOKEN", "")),
		Secure: env.Bool("STORAGE_S3_SECURE", true),
		Region: env.Get("STORAGE_S3_REGION", ""),
	})
	if err != nil {
		return err
	}
	driver.client = client
	driver.bucket = env.Get("STORAGE_S3_BUCKET", "pasty")
	return nil
}

// Terminate terminates the AWS S3 storage driver (does nothing, because the AWS S3 storage driver does not need any termination)
func (driver *S3Driver) Terminate() error {
	return nil
}

// ListIDs returns a list of all existing paste IDs
func (driver *S3Driver) ListIDs() ([]string, error) {
	// Define the IDs slice
	var ids []string

	// Fill the IDs slice
	channel := driver.client.ListObjects(context.Background(), driver.bucket, minio.ListObjectsOptions{})
	for object := range channel {
		if object.Err != nil {
			return nil, object.Err
		}
		ids = append(ids, strings.TrimSuffix(object.Key, ".json"))
	}

	// Return the IDs slice
	return ids, nil
}

// Get loads a paste
func (driver *S3Driver) Get(id string) (*pastes.Paste, error) {
	// Read the object
	object, err := driver.client.GetObject(context.Background(), driver.bucket, id+".json", minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}

	// Unmarshal the object into a paste
	paste := new(pastes.Paste)
	err = json.Unmarshal(data, &paste)
	if err != nil {
		return nil, err
	}
	return paste, nil
}

// Save saves a paste
func (driver *S3Driver) Save(paste *pastes.Paste) error {
	// Marshal the paste
	jsonBytes, err := json.Marshal(paste)
	if err != nil {
		return err
	}

	// Put the object
	reader := bytes.NewReader(jsonBytes)
	_, err = driver.client.PutObject(context.Background(), driver.bucket, paste.ID+".json", reader, reader.Size(), minio.PutObjectOptions{
		ContentType: "application/json",
	})
	return err
}

// Delete deletes a paste
func (driver *S3Driver) Delete(id string) error {
	return driver.client.RemoveObject(context.Background(), driver.bucket, id+".json", minio.RemoveObjectOptions{})
}
