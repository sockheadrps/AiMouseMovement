package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	mongoClient *mongo.Client
}

func NewClient() Client {
	return Client{}
}

func (c Client) BuildMongoOptions() *options.ClientOptions {
	// set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// pls dont bomb my free tier mongo db
	return options.Client().ApplyURI("mongodb+srv://sockheadrps:s3eUDQdQqO82UYH2@cluster0.g5abd.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
}

func (c Client) Connect(ctx context.Context, opts *options.ClientOptions) (error) {
	// Create a new client and connect to the server
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("error while connecting to Mongo client: %w", err)
	}

	c.mongoClient = mongoClient

	return nil
}

func (c Client) Disconnect(ctx context.Context) error {
	err := c.mongoClient.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("error while disconnecting from Mongo client: %w", err)
	}

	return nil
}

func (c Client) Run(ctx context.Context, db string, cmd interface{}) (bson.M, error) {
	var result bson.M
	if err := c.mongoClient.Database(db).RunCommand(ctx, cmd).Decode(&result); err != nil {
		return result, fmt.Errorf("error while running command: %w", err)
	}
	return result, nil
}

func (c Client) Collection(database, collection string) *mongo.Collection {
	return c.mongoClient.Database(database).Collection(collection)
}