package handlers

import (
	"context"
	"dbServer/api"
	"dbServer/db"
	"dbServer/models"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDisasterList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	disasters, err := api.GetDisaster(ctx)
	if err != nil {
		fmt.Println("Error getting disaster: ", err)
		c.JSON(500, gin.H{"message": "Error getting disaster"})
		return
	}

	log.Printf("Disasters: %v", disasters)

	people, err := parseResponse(c, disasters)
	if err != nil {
		fmt.Println("Error parsing response: ", err)
		c.JSON(500, gin.H{"message": "Error parsing response"})
		return
	}

	log.Printf("People: %v", people)

	/*
		for _, person := range people {
			_, err := db.Mongo.Collections["Message"].InsertOne(context.TODO(), person)
			if err != nil {
				fmt.Println("Error inserting person: ", err)
				c.JSON(500, gin.H{"message": "Error inserting person"})
				return
			}
		}
	*/

	result, err := db.Mongo.Collections["Message"].InsertMany(context.TODO(), people)
	if err != nil {
		fmt.Println("Error inserting people: ", err)
		c.JSON(500, gin.H{"message": "Error inserting people: " + err.Error()})
		return
	}
	resultString := fmt.Sprintf("Inserted %d people", len(result.InsertedIDs))
	fmt.Println(resultString)

	c.JSON(200, resultString)

	//fmt.Println("Disasters: ", disasters)
	//c.JSON(200, disasters)

	return
}

func isMissingPersonMsg(msg string) (bool, error) {
	result := strings.Contains(msg, "찾습니다") &&
		(strings.Contains(msg, "씨(남,") || strings.Contains(msg, "씨(여,"))

	if result == false {
		return false, fmt.Errorf("Not a missing person message: msg: %s", msg)
	}
	return result, nil
}

func ParseDisaster(disaster models.Disaster) (*models.Person, error) {
	sn := strconv.Itoa(disaster.SN)

	timestamp, err := time.Parse("2006/01/02 15:04:05", disaster.Timestamp)
	if err != nil {
		fmt.Println("Error parsing timestamp: ", err)
		return nil, err
	}
	dateTime := primitive.NewDateTimeFromTime(timestamp)

	regions := strings.Split(disaster.Region, ",")

	missing, err := ParseMessage(disaster.Message)
	if err != nil {
		fmt.Println("Error parsing message: ", err)
		return nil, err
	}

	person := &models.Person{
		SN:        sn,
		Timestamp: dateTime,
		Region:    regions,
		Missing:   *missing,
		Accessed:  false,
	}
	return person, nil
}

func ParseMessage(message string) (*models.Missing, error) {
	nameRegex := regexp.MustCompile(`(\S+)씨\((\S+),\s*(\d+)세\)`)
	matches := nameRegex.FindStringSubmatch(message)

	if len(matches) != 4 {
		return nil, fmt.Errorf("이름, 성별, 나이 파싱 실패")
	}

	age, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("나이 변환 실패: %v", err)
	}

	descRegex := regexp.MustCompile(`를 찾습니다-(.+?)(?:\r\n|\s+)vo\.la`)
	descMatches := descRegex.FindStringSubmatch(message)

	if len(descMatches) != 2 {
		return nil, fmt.Errorf("설명 파싱 실패")
	}

	missing := &models.Missing{
		Name:        matches[1],
		Age:         age,
		Sex:         matches[2],
		Description: descMatches[1],
	}

	return missing, nil
}

func parseResponse(c *gin.Context, disasters []models.Disaster) ([]interface{}, error) {
	var people []interface{}

	for _, disaster := range disasters {
		if isMiss, err := isMissingPersonMsg(disaster.Message); isMiss == false {
			log.Println(err)
			continue
		}

		if alreadyExists(strconv.Itoa(disaster.SN)) {
			log.Printf("Person already exists: %d", disaster.SN)
			return people, nil
		}
		person, err := ParseDisaster(disaster)
		if err != nil {
			fmt.Printf("Error parsing disaster(sn: %d): %v", disaster.SN, err)
			continue
		}

		people = append(people, *person)
		log.Println("Person: ", person)
	}

	return people, nil
}

func alreadyExists(sn string) bool {
	existPerson := db.Mongo.Collections["Message"].FindOne(context.TODO(), bson.M{"sn": sn})
	if existPerson.Err() != nil {
		return false
	}
	return true
}
