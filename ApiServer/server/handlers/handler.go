package handlers

import (
	"apiServer/db"
	"apiServer/models"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewDataExist(c *gin.Context) {
	collection := db.Mongo.Collections["Message"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	unaccessed := collection.FindOne(ctx, bson.M{"accessed": false})

	if unaccessed != nil {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Data exist",
		})
		return
	}
}

func SearchBySN(c *gin.Context) {
	collection := db.Mongo.Collections["Message"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	missingPerson := models.Person{}
	sn := c.Param("sn")
	//log.Fatalf("sn: %v", sn)

	err := collection.FindOne(ctx, bson.M{"sn": sn}).Decode(&missingPerson)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Data not found",
			"error":   err,
		})
		return
	}
	fmt.Printf("missingPerson: %v\n", missingPerson)

	c.JSON(200, gin.H{
		"sn":        missingPerson.SN,
		"timestamp": missingPerson.Timestamp,
		"region":    missingPerson.Region,
		"missing":   missingPerson.Missing,
		"accessed":  missingPerson.Accessed,
	})

}

func GetAdditionalList(c *gin.Context) {
	collection := db.Mongo.Collections["Message"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var people []models.Person

	cursor, err := collection.Find(ctx, bson.M{"accessed": false})
	if err != nil {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Data not found",
		})
		return
	}

	if err := cursor.All(c, &people); err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"message": "Internal server error",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": people,
	})

}

func GetAllList(c *gin.Context) {
	collection := db.Mongo.Collections["Message"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var people []models.Person

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Data not found",
		})
		return
	}

	if err := cursor.All(c, &people); err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"message": "Internal server error",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": people,
	})

}

func SearchByRegion(c *gin.Context) {
	collection := db.Mongo.Collections["Message"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var people []models.Person
	region := c.Param("region")

	filter := bson.M{
		"region": bson.M{
			"$elemMatch": bson.M{
				"$regex": primitive.Regex{
					Pattern: region,
					Options: "i",
				},
			},
		},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Data not found",
			"error":   err,
		})
		return
	}

	if err := cursor.All(c, &people); err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"message": "Internal server error",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": people,
	})

}
