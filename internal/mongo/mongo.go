package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDBClient initializes a new MongoDB client
func NewMongoClient(ctx context.Context, uri string, dbName string) (*MongoClient, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB")

	// Return the MongoDB client instance
	return &MongoClient{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// InsertFile stores a file's metadata and content in MongoDB
func (m *MongoClient) InsertFile(ctx context.Context, collectionName string, fileData map[string]interface{}) error {
	collection := m.database.Collection(collectionName)

	_, err := collection.InsertOne(ctx, fileData)
	if err != nil {
		return fmt.Errorf("failed to insert file into MongoDB: %w", err)
	}

	fmt.Println("File successfully inserted into MongoDB")
	return nil
}
