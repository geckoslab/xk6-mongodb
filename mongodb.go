package xk6_mongodb

import (
	"context"
	k6Modules "go.k6.io/k6/js/modules"
	mongoPrimitive "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOpts "go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func init() {
	k6Modules.Register("k6/x/mongodb", new(MongoDb))
}

type MongoDb struct{}

type Connection struct {
	Client *mongo.Client
}

func (*MongoDb) Connect(url string) *Connection {
	connectionOpts := mongoOpts.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), connectionOpts)

	if err != nil {
		log.Fatalf("Error when establishing connection to MongoDB: %v", err)
	}

	return &Connection{
		Client: client,
	}
}

func (*MongoDb) NewId() mongoPrimitive.ObjectID {
	return mongoPrimitive.NewObjectID()
}

func (*MongoDb) NewIdFromHex(hex string) mongoPrimitive.ObjectID {
	objectId, err := mongoPrimitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Fatalf("Error when creating ObjectID from hex: %v", err)
	}

	return objectId
}

func (connection *Connection) Close() {
	err := connection.Client.Disconnect(context.Background())
	if err != nil {
		log.Fatalf("Error when closing connection to MongoDB: %v", err)
	}
}

func (connection *Connection) Insert(dbName string, collName string, doc interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	_, err := collection.InsertOne(context.Background(), doc)

	if err != nil {
		log.Fatalf("Error when inserting document to MongoDB: %v", err)
		return err
	}

	return nil
}

func (connection *Connection) InsertMany(dbName string, collName string, docs []interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)
	_, err := collection.InsertMany(context.Background(), docs)

	if err != nil {
		log.Fatalf("Error when inserting documents to MongoDB: %v", err)
		return err
	}

	return nil
}
