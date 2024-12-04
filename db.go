package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

func NewMongoClient(uri, database, collection string) (*MongoClient, error) {
	options := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	newClient := MongoClient{
		Client:     client,
		Database:   database,
		Collection: collection,
	}

	return &newClient, nil

}

// Disconnects will disconnect from the Mongo database
func (m *MongoClient) Disconnect() (bool, error) {
	if err := m.Client.Disconnect(context.TODO()); err != nil {
		return false, err
	}
	return true, nil
}
