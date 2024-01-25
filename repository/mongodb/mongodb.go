package mongodb

import (
	"api/config"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoStore struct {
	client *mongo.Client
	dbName string
}

func NewMongoConnection(connectURI, databaseName string) (*MongoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(connectURI)
	opts.SetAppName(config.ServiceName)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = createIndexes(ctx, client, databaseName)
	if err != nil {
		return nil, err
	}

	return &MongoStore{client: client, dbName: databaseName}, nil
}

func createIndexes(ctx context.Context, client *mongo.Client, dbName string) error {
	// Defining an index model for a 2dsphere index on the location field.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"location", "2dsphere"}},
		Options: options.Index().SetName("location_2dsphere"),
	}

	_, err := client.Database(dbName).Collection("users").Indexes().CreateOne(ctx, indexModel)
	return err
}

// coll is a helper method to get a collection from the MongoDB store.
func (conn *MongoStore) coll(name string) *mongo.Collection {
	return conn.client.Database(conn.dbName).Collection(name)
}
