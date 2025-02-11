package xk6_mongodb

import (
	"context"
	"log/slog"

	k6Modules "go.k6.io/k6/js/modules"
	mongoBson "go.mongodb.org/mongo-driver/bson"
	mongoPrimitive "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOpts "go.mongodb.org/mongo-driver/mongo/options"
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
		slog.Error("Error when establishing connection to MongoDB.", "error", err)
	}

	return &Connection{
		Client: client,
	}
}

func (connection *Connection) Close() error {
	err := connection.Client.Disconnect(context.Background())
	if err != nil {
		slog.Error("Error when closing connection to MongoDB.", "error", err)
		return err
	}

	return nil
}

func (*MongoDb) NewObjectID() mongoPrimitive.ObjectID {
	return mongoPrimitive.NewObjectID()
}

func (*MongoDb) ObjectIDFromHex(hex string) mongoPrimitive.ObjectID {
	objectId, err := mongoPrimitive.ObjectIDFromHex(hex)
	if err != nil {
		slog.Error("Error when creating ObjectID from hex.", "error", err)
	}

	return objectId
}

func (connection *Connection) Insert(dbName string, collName string, doc interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	_, err := collection.InsertOne(context.Background(), doc)

	if err != nil {
		slog.Error("Error when inserting document to MongoDB.", "error", err)
	}

	return nil
}

func (connection *Connection) InsertMany(dbName string, collName string, docs []interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)
	_, err := collection.InsertMany(context.Background(), docs)

	if err != nil {
		slog.Error("Error when inserting documents to MongoDB.", "error", err)
	}

	return nil
}

func (connection *Connection) Upsert(dbName string, collName string, filter interface{}, update interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	_, err := collection.UpdateOne(context.Background(), filter, update, mongoOpts.Update().SetUpsert(true))
	if err != nil {
		slog.Error("Error when upserting document in MongoDB.", "error", err)
	}

	return nil
}

func (connection *Connection) FindOne(dbName string, collName string, filter interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	var result mongoBson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		slog.Error("Error when finding document in MongoDB.", "error", err)
	}

	return nil
}

func (connection *Connection) Find(dbName string, collName string, filter interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		slog.Error("Error when finding documents in MongoDB.", "error", err)
	}

	var results []mongoBson.M
	if err = cur.All(context.Background(), &results); err != nil {
		slog.Error("Error while decoding documents.", "error", err)
	}

	return nil
}

func (connection *Connection) FindAll(dbName string, collName string) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	cur, err := collection.Find(context.Background(), mongoBson.M{})
	if err != nil {
		slog.Error("Error when finding documents in MongoDB.", "error", err)
	}

	var results []mongoBson.M
	if err = cur.All(context.Background(), &results); err != nil {
		slog.Error("Error while decoding documents.", "error", err)
	}

	return nil
}

func (connection *Connection) UpdateOne(dbName string, collName string, filter interface{}, update interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		slog.Error("Error when updating document in MongoDB.", "error", err)
	}

	return nil
}

func (connection *Connection) UpdateMany(dbName string, collName string, filter interface{}, update interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	_, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		slog.Error("Error when updating documents in MongoDB.", "error", err)
	}

	return nil
}

func (connection *Connection) DeleteOne(dbName string, collName string, filter interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		slog.Error("Error when deleting document in MongoDB.", "error", err)
	}

	return nil
}

func (connection *Connection) DeleteMany(dbName string, collName string, filter interface{}) error {
	collection := connection.Client.Database(dbName).Collection(collName)

	_, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		slog.Error("Error when deleting documents in MongoDB.", "error", err)
	}

	return nil
}
