package upload_svc

import (
	"context"
	"errors"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

type DBStruct struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type DB interface {
	Insert(info ImgInfo) error
	Query(imgID string) (ImgInfo, error)
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

func storeImgInfo(info ImgInfo) error {
	return db.Insert(info)
}

func getImgInfo(imgID string) (ImgInfo, error) {
	return db.Query(imgID)
}

func addPerData() error {
	return db.(DBStruct).preData()
}

func (dbs DBStruct) preData() error {
	predata := ImgInfo{
		ImgID:  "test",
		ImgURL: "testURL",
		User:   "testUser",
	}
	_, err := dbs.Collection.InsertOne(context.TODO(), predata)
	return err
}

func (dbs DBStruct) Insert(info ImgInfo) error {
	if !info.Valid() {
		return errors.New("info not valid")
	}
	_, err := dbs.Collection.InsertOne(context.TODO(), info)
	return err
}

func (dbs DBStruct) Query(imgID string) (ImgInfo, error) {
	var result ImgInfo
	err := dbs.Collection.FindOne(context.TODO(), bson.M{
		"imgid": imgID,
	}).Decode(&result)

	return result, err
}
