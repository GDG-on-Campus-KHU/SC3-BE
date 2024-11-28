package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	SN        string             `bson:"sn"`
	Timestamp primitive.DateTime `bson:"timestamp"`
	Region    []string           `bson:"region"`
	Missing   Missing            `bson:"missing"`
	Accessed  bool               `bson:"accessed"`
}

type Missing struct {
	Name        string `bson:"name"`
	Age         int    `bson:"age"`
	Sex         string `bson:"sex"`
	Description string `bson:"description"`
}
