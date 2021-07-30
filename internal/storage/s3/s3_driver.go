package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/shared"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Driver represents the AWS S3 storage driver
type S3Driver struct {
	client *minio.Client
	bucket string
}

// Initialize initializes the AWS S3 storage driver
func (driver *S3Driver) Initialize() error {
	client, err := minio.New(config.Current.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Current.S3.AccessKeyID, config.Current.S3.SecretAccessKey, config.Current.S3.SecretToken),
		Secure: config.Current.S3.Secure,
		Region: config.Current.S3.Region,
	})
	if err != nil {
		return err
	}
	driver.client = client
	driver.bucket = config.Current.S3.Bucket
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
func (driver *S3Driver) Get(id string) (*shared.Paste, error) {
	// Read the object
	object, err := driver.client.GetObject(context.Background(), driver.bucket, id+".json", minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(object)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal the object into a paste
	paste := new(shared.Paste)
	err = json.Unmarshal(data, &paste)
	if err != nil {
		return nil, err
	}
	return paste, nil
}

// Save saves a paste
func (driver *S3Driver) Save(paste *shared.Paste) error {
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

// Cleanup cleans up the expired pastes
func (driver *S3Driver) Cleanup() (int, error) {
	// Retrieve all paste IDs
	ids, err := driver.ListIDs()
	if err != nil {
		return 0, err
	}

	// Define the amount of deleted items
	deleted := 0

	// Loop through all pastes
	for _, id := range ids {
		// Retrieve the paste object
		paste, err := driver.Get(id)
		if err != nil {
			return 0, err
		}

		// Delete the paste if it is expired
		lifetime := config.Current.AutoDelete.Lifetime
		if paste.Created+int64(lifetime.Seconds()) < time.Now().Unix() {
			err = driver.Delete(id)
			if err != nil {
				return 0, err
			}
			deleted++
		}
	}
	return deleted, nil
}
