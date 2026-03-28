package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var mongoClient *mongo.Client
var dbName string
var UserCollectionName string = "Users"
var NoteCollectionName string = "Notes"

func ConnectDB(mongoURI, database string) error {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	mongoClient = client
	dbName = database

	return nil
}

func DisconnectDB() error {
	if mongoClient == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	log.Println("✓ MongoDB connection closed")

	return nil
}

func GetDB() *mongo.Database {
	if mongoClient == nil {
		panic("MongoDB client is not initialized. Call ConnectDB() first")
	}
	return mongoClient.Database(dbName)
}
