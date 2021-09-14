package account_svc

import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-kit/log"
)

type DBStruct struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type DB interface {
	UserLogin(UserInfo) (string, UserInfo)
	CreateUser(UserInfo) bool
}

var db DB
var Logger log.Logger = init_logger()

func initDB() {
	client := connectDB()
	_db := DBStruct{
		Client:     client,
		Collection: client.Database("auth").Collection("auth_info"),
	}
	db = _db
}

func init_logger() log.Logger {
	return log.NewLogfmtLogger(log.StdlibWriter{})
}

func connectDB() *mongo.Client {
	var mongoOnce sync.Once
	var client *mongo.Client
	var clientInstanceError error
	var CONNECTIONSTRING = "mongodb://" + os.Getenv("DB_URL")

	Logger.Log("status", "create database connection with : "+CONNECTIONSTRING)
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

func (dbs DBStruct) UserLogin(user UserInfo) (string, UserInfo) {
	Logger.Log("status", "user login")
	if dbs.tryLogin(user) {
		user = dbs.getUserInfo(user)
		return generateJWTToken(user), user
	} else {
		return "", UserInfo{}
	}
}

func (dbs DBStruct) tryLogin(user UserInfo) bool {
	var result UserInfo
	err := dbs.Collection.FindOne(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email, "password": user.Password},
			bson.M{"name": user.Name, "password": user.Password},
		},
	}).Decode(&result)

	Logger.Log("try login", err == nil)
	return err == nil
}

func (dbs DBStruct) getUserInfo(user UserInfo) UserInfo {
	var result UserInfo
	err := dbs.Collection.FindOne(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email, "password": user.Password},
			bson.M{"name": user.Name, "password": user.Password},
		},
	}).Decode(&result)

	Logger.Log("GetUserInfo", err == nil)
	return result
}

func (dbs DBStruct) CreateUser(user UserInfo) bool {
	Logger.Log("status", "create user")
	if !dbs.valid(user) {
		Logger.Log("create user", false, "msg", "input not valid, probably miss email, name or password")
		return false
	}
	if dbs.userExisted(user) {
		Logger.Log("create user", false)
		return false
	}
	toInsert := UserInfo{
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}
	_, err := dbs.Collection.InsertOne(context.TODO(), toInsert)
	Logger.Log("create user", err == nil)
	return err == nil
}

func (dbs DBStruct) valid(user UserInfo) bool {
	return user.Email != "" &&
		user.Name != "" &&
		user.Password != ""
}

func (dbs DBStruct) userExisted(user UserInfo) bool {
	var result UserInfo
	err := dbs.Collection.FindOne(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"name": user.Name},
		},
	}).Decode(&result)
	return err != mongo.ErrNoDocuments
}
