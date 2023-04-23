package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client represents a MongoDB database
type Client struct {
	mongoClient *mongo.Client
	database    *mongo.Database
	collection  *mongo.Collection
}

// Init creates a new MongoDB database
func (c *Client) Init(ctx context.Context, db, collection string) error {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	database := client.Database(db)
	col := database.Collection(collection)

	c.mongoClient = client
	c.database = database
	c.collection = col
	return nil
}

// Close closes the database connection
func (c *Client) Close(ctx context.Context) error {
	err := c.mongoClient.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %v", err)
	}

	return nil
}
