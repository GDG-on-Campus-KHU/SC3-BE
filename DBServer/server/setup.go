package server

import (
	"context"
	"dbServer/api"
	"dbServer/db"
	"dbServer/server/handlers"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	Router *gin.Engine
}

var (
	Server *GinServer = &GinServer{}
)

func New() {
	Server.Router = gin.Default()

	Server.Router.GET("/api/disaster/list", handlers.GetDisasterList)
	go startPeriodicFetch()
}

func startPeriodicFetch() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Starting periodic disaster list fetch...")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			disasters, err := api.GetDisaster(ctx, "1", "300", time.Now().AddDate(0, 0, -1).Format("20060102"))
			if err != nil {
				log.Printf("Error fetching disasters: %v", err)
				cancel()
				continue
			}

			people, err := handlers.ParseResponse(disasters)
			if err != nil {
				log.Printf("Error parsing response: %v", err)
				cancel()
				continue
			}

			if len(people) > 0 {
				result, err := db.Mongo.Collections["Message"].InsertMany(ctx, people)
				if err != nil {
					log.Printf("Error inserting people: %v", err)
				} else {
					log.Printf("Inserted %d people", len(result.InsertedIDs))
				}
			} else {
				log.Println("No new missing person found")
			}

			cancel()
			log.Println("Completed periodic disaster list fetch")
		}
	}
}
