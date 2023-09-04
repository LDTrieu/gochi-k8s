package data

import (
	"context"
	"errors"
	"log-service/utils"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	logDBClient *mongo.Client
	logDB       *mongo.Database
	//logColl     *mongo.Collection
)

const (
	mongodbURL = "MONGODB_URL"
	userName   = "MONGODB_USER"
	password   = "MONGODB_PASSWORD"
	logDBName  = "DB_LOG"
)

type LogModel struct {
	ID        string
	Value     string
	CreatedAt time.Time
}

func InitDB() error {

	url, err := utils.ViperEnvVariable(mongodbURL)
	if err != nil {
		return errors.New("not find database url")
	}

	dbName, err := utils.ViperEnvVariable(logDBName)
	if err != nil {
		return errors.New("not find database name")
	}

	clientOptions := options.Client().ApplyURI(url)
	logDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	err = logDBClient.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	logDB = logDBClient.Database(dbName)
	//logColl = logDB.Collection("logs")

	return nil
}

func runDBTask(task func(db *mongo.Database) error) error {
	if logDB == nil {
		err := InitDB()
		if err != nil {
			return err
		}
	}
	return task(logDB)
}

func InsertLog(ctx context.Context, value string) error {
	task := func(db *mongo.Database) error {
		log := &LogModel{
			Value:     value,
			CreatedAt: time.Now(),
		}

		_, err := db.Collection("log", nil).InsertOne(ctx, &log)
		if err != nil {
			return err
		}
		return err
	}

	return runDBTask(task)
}
