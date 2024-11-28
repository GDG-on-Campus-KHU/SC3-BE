package db

import (
	cfg "apiServer/config"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client      *mongo.Client
	Collections map[string]*mongo.Collection
}

var (
	Mongo *DB = &DB{}
)

func New(c context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.AppConfig.DB.URL).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database(cfg.AppConfig.DB.Database).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	Mongo.Client = client
	Mongo.Collections = make(map[string]*mongo.Collection)

	database := client.Database(cfg.AppConfig.DB.Database)
	Mongo.Collections["Message"] = database.Collection("Message")

	return ctx, cancel
}

func DisconnectDB() {
	if err := Mongo.Client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
