package xk6_mongodb

import (
	"context"
	"log"
	"reflect"

	k6Modules "go.k6.io/k6/js/modules"
	mongoPrimitives "go.mongodb.org/mongo-driver/bson/primitive"
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

func (m *MongoDb) Connect(url string) *Connection {
	connectionOpts := mongoOpts.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), connectionOpts)

	if err != nil {
		log.Fatalf("Error when establishing connection to MongoDB: %v", err)
	}

	return &Connection{
		Client: client,
	}
}

func (connection *Connection) TransformDoc(doc interface{}) interface{} {
	// Transform all $oid (nested) to ObjectID
	for key, value := range doc.(map[string]interface{}) {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			if _, ok := value.(map[string]interface{})["$oid"]; ok {
				newObjectID, err := mongoPrimitives.ObjectIDFromHex(value.(map[string]interface{})["$oid"].(string))
				if err != nil {
					log.Fatalf("Error when converting string to ObjectID: %v", err)
				}

				doc.(map[string]interface{})[key] = newObjectID
			} else {
				// TODO: Handle nested object
			}
		}
	}

	return doc
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
