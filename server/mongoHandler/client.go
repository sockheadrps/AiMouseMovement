package mongoHandler

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	mongoClient *mongo.Client
}

func NewClient() Client {
	return Client{}
}

func (c *Client) BuildMongoOptions(mongoUrl string) *options.ClientOptions {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	return options.Client().ApplyURI(mongoUrl).SetServerAPIOptions(serverAPI)
}

func (c *Client) Connect(ctx context.Context, opts *options.ClientOptions) error {
	// Create a new client and connect to the server
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("error while connecting to Mongo client: %w", err)
	}
	c.mongoClient = mongoClient

	return nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	err := c.mongoClient.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("error while disconnecting from Mongo client: %w", err)
	}

	return nil
}

func (c *Client) Run(ctx context.Context, db string, cmd interface{}) (bson.M, error) {
	var result bson.M
	if err := c.mongoClient.Database(db).RunCommand(ctx, cmd).Decode(&result); err != nil {
		return result, fmt.Errorf("error while running command: %w", err)
	}
	return result, nil
}

func (c *Client) Collection(database, collection string) *mongo.Collection {
	return c.mongoClient.Database(database).Collection(collection)
}

func (c *Client) Insert(ctx context.Context, database, collection string, document interface{}) error {
	_, err := c.mongoClient.Database(database).Collection(collection).InsertOne(ctx, document, nil)
	if err != nil {
		return fmt.Errorf("error while inserting document: %w", err)
	}
	return nil
}

func (c *Client) RemoveByID(ctx context.Context, database, collection, documentID string) (int64, error) {
	idPrimitive, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
		fmt.Println("primitive.ObjectIDFromHex ERROR")
	}

	filter := bson.M{"_id": idPrimitive}

	result, err := c.mongoClient.Database(database).Collection(collection).DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("primitive. ERROR")
		return 0, fmt.Errorf("error while removing document: %w", err)
	}

	// Check if a document was deleted
	fmt.Println(result)
	deletedCount := result.DeletedCount
	return deletedCount, nil
}
