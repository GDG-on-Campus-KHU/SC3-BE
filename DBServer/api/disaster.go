package api

import (
	"context"
	"dbServer/models"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type APIResponse struct {
	Header struct {
		ResultMsg  string      `json:"resultMsg"`
		ResultCode string      `json:"resultCode"`
		ErrorMsg   interface{} `json:"errorMsg"`
	} `json:"header"`
	NumOfRows  int               `json:"numOfRows"`
	PageNo     int               `json:"pageNo"`
	TotalCount int               `json:"totalCount"`
	Body       []models.Disaster `json:"body"`
}

func GetDisaster(c context.Context, pageNum, rowNum, crtDt string) ([]models.Disaster, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	serviceKey := os.Getenv("API_SERVICE_KEY")
	uri := os.Getenv("API_URI")

	url := "https://www.safetydata.go.kr/" + uri + "?serviceKey=" + serviceKey
	url = url + "&numOfRows=" + rowNum
	url = url + "&pageNo=" + pageNum
	url = url + "&crtDt=" + crtDt
	url = url + "&type=json"
	log.Printf("URL: %v", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse APIResponse
	//responseBody := []models.Disaster{}
	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}
	return apiResponse.Body, nil
}
