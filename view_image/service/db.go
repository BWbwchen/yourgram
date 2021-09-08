package service

import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBStruct struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

var db DB

func initDB() {
	client := connectDB()
	_db := DBStruct{
		Client:     client,
		Collection: client.Database("img").Collection("img_info"),
	}
	db = _db
}
func connectDB() *mongo.Client {
	var mongoOnce sync.Once
	var client *mongo.Client
	var clientInstanceError error
	var CONNECTIONSTRING = "mongodb://" + os.Getenv("DB_URL")

	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		_client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = _client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		client = _client
	})
	if clientInstanceError != nil {
		panic(clientInstanceError)
	}
	return client
}

type Result struct {
	URL string `json:"imgid"`
}

func (dbs DBStruct) Query(UserID string) ([]string, error) {
	var urls []string
	cursor, err := dbs.Collection.Find(context.TODO(), bson.M{
		"user": UserID,
	})

	if err != nil {
		return urls, err
	}

	var results []bson.M
	err = cursor.All(context.TODO(), &results)

	var imgid string
	for _, result := range results {
		// TODO : eliminate this environment variable
		imgid = os.Getenv("baseurl") + "/" + result["imgid"].(string)
		urls = append(urls, imgid)
	}

	return urls, err
}
