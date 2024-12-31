package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	YoutubeApiKey1     string
	YoutubeApiKey2     string
	YoutubeApiKey3     string
	PostgresServiceURI string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Println("Error loading config/.env\nerror:", err.Error())
		// Still Continue
	}

	AppConfig = Config{
		YoutubeApiKey1:     getEnv("YOUTUBE_API_KEY1", "NA"),
		YoutubeApiKey2:     getEnv("YOUTUBE_API_KEY2", ""),
		YoutubeApiKey3:     getEnv("YOUTUBE_API_KEY3", ""),
		PostgresServiceURI: getEnv("POSTGRES_SERVICE_URI", "NA"),
	}
	log.Println("Loaded Configurations...")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
