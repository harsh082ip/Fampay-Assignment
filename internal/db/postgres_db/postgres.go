package postgres_db

import (
	"log"

	"github.com/harsh082ip/Fampay-Assignment/internal/config"
	"github.com/harsh082ip/Fampay-Assignment/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.AppConfig.PostgresServiceURI
	if dsn == "" {
		log.Fatal("Database Connection String not provided")
	}

	// Connect to postgres
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
	}

	if err := DB.AutoMigrate(&models.Video{}); err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}
}
