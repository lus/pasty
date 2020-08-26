package storage

import (
	"context"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// MongoDBDriver represents the MongoDB storage driver
type MongoDBDriver struct {
	client     *mongo.Client
	database   string
	collection string
}

// Initialize initializes the MongoDB storage driver
func (driver *MongoDBDriver) Initialize() error {
	// Define the context for the following database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB host
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.Get("STORAGE_MONGODB_CONNECTION_STRING", "mongodb://pasty:pasty@example.host/pasty")))
	if err != nil {
		return err
	}

	// Ping the MongoDB host
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	// Set the driver attributes
	driver.client = client
	driver.database = env.Get("STORAGE_MONGODB_DATABASE", "pasty")
	driver.collection = env.Get("STORAGE_MONGODB_COLLECTION", "pastes")
	return nil
}

// Terminate terminates the MongoDB storage driver
func (driver *MongoDBDriver) Terminate() error {
	return driver.client.Disconnect(context.TODO())
}

// ListIDs returns a list of all existing paste IDs
func (driver *MongoDBDriver) ListIDs() ([]string, error) {
	// Define the collection to use for this database operation
	collection := driver.client.Database(driver.database).Collection(driver.collection)

	// Define the context for the following database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve all paste documents
	result, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Decode all paste documents
	var pasteSlice []pastes.Paste
	err = result.All(ctx, &pasteSlice)
	if err != nil {
		return nil, err
	}

	// Read and return the IDs of all paste objects
	var ids []string
	for _, paste := range pasteSlice {
		ids = append(ids, paste.ID)
	}
	return ids, nil
}

// Get loads a paste
func (driver *MongoDBDriver) Get(id string) (*pastes.Paste, error) {
	// Define the collection to use for this database operation
	collection := driver.client.Database(driver.database).Collection(driver.collection)

	// Define the context for the following database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to retrieve the corresponding paste document
	filter := bson.M{"_id": id}
	result := collection.FindOne(ctx, filter)
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// Return the retrieved paste object
	paste := new(pastes.Paste)
	err = result.Decode(paste)
	return paste, err
}

// Save saves a paste
func (driver *MongoDBDriver) Save(paste *pastes.Paste) error {
	// Define the collection to use for this database operation
	collection := driver.client.Database(driver.database).Collection(driver.collection)

	// Define the context for the following database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the paste object
	_, err := collection.InsertOne(ctx, paste)
	return err
}

// Delete deletes a paste
func (driver *MongoDBDriver) Delete(id string) error {
	// Define the collection to use for this database operation
	collection := driver.client.Database(driver.database).Collection(driver.collection)

	// Define the context for the following database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete the document
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}
