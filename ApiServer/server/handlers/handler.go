package handlers

import (
	"apiServer/db"
	"apiServer/models"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
