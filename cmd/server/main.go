package main

import (
	"log"

	"github.com/harsh082ip/Fampay-Assignment/internal/config"
	"github.com/harsh082ip/Fampay-Assignment/internal/db/postgres_db"
	"github.com/harsh082ip/Fampay-Assignment/internal/router"
	"github.com/harsh082ip/Fampay-Assignment/internal/videos"
)

const (
	WEBPORT = ":8000"
)

func main() {

	config.LoadConfig()
	postgres_db.InitDB()
	go videos.FetchYouTubeVideos("cricket")
	r := router.SetupRouter()
	if err := r.Run(WEBPORT); err != nil {
		log.Fatal("Error Starting Server on ", WEBPORT)
	}
}
