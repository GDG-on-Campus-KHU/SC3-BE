package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB  MongoConfig
	Srv ServerConfig
}

type MongoConfig struct {
	URL      string
	Database string
}

type ServerConfig struct {
	Port string `default:"5000" envconfig:"SERVER_PORT"`
}

var (
	AppConfig *Config = &Config{}
)

func DBInit() MongoConfig {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PWD")
	dbUri := os.Getenv("DB_URI")
	url := "mongodb+srv://" + dbUser + ":" + dbPassword + dbUri
	mongo := MongoConfig{
		URL:      url,
		Database: "DisasterText",
	}
	return mongo
}

func ServerInit() ServerConfig {
	server := ServerConfig{
		Port: os.Getenv("SERVER_PORT"),
	}
	return server
}

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	AppConfig.DB = DBInit()
	AppConfig.Srv = ServerInit()
}
