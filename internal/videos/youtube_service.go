package videos

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/harsh082ip/Fampay-Assignment/internal/config"
	"github.com/harsh082ip/Fampay-Assignment/internal/db/postgres_db"
	"github.com/harsh082ip/Fampay-Assignment/internal/models"
)

func FetchYouTubeVideos(searchQuery string) {
	log.Println("Fetching Started...")
	baseURL := "https://www.googleapis.com/youtube/v3/search"
	nextPageToken := ""
	ytApiKey := config.AppConfig.YoutubeApiKey
	var videosToInsert []models.Video
	for {
		url := fmt.Sprintf(
			"%s?part=snippet&maxResults=10&q=%s&type=video&pageToken=%s&key=%s",
			baseURL, searchQuery, nextPageToken, ytApiKey,
		)

		// make HTTP request
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error Fetching Videos Info, error: %v", err)
		}

		// Parse body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response, error: %v", err)
		}

		var apiResponse models.YoutubeApiResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			log.Fatalf("Failed to parse JSON response: %v", err)
		}

		for _, item := range apiResponse.Items {
			publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishTime)
			if err != nil {
				log.Printf("Failed to parse publishedAt for video %s: %v", item.ID.VideoID, err)
				continue
			}

			// Marshal thumbnails into JSON
			thumbnailURLs, err := json.Marshal([]string{
				item.Snippet.Thumbnails.Default.URL,
				item.Snippet.Thumbnails.Medium.URL,
				item.Snippet.Thumbnails.High.URL,
			})
			if err != nil {
				log.Printf("Failed to marshal thumbnails for video %s: %v", item.ID.VideoID, err)
				continue
			}

			// Handle nullable fields
			var description *string
			if item.Snippet.Description != "" {
				description = &item.Snippet.Description
			}

			var thumbnails *string
			if len(thumbnailURLs) > 0 {
				thumbnails = new(string)
				*thumbnails = string(thumbnailURLs)
			}

			video := models.Video{
				VideoID:       item.ID.VideoID,
				Title:         item.Snippet.Title,
				Description:   description,
				PublishedAt:   publishedAt,
				ThumbnailURLs: thumbnails,
			}

			videosToInsert = append(videosToInsert, video)
		}

		if len(videosToInsert) > 0 {
			result := postgres_db.DB.CreateInBatches(videosToInsert, 100)
			if result.Error != nil {
				log.Printf("Error inserting videos into the database: %v", result.Error)
			}

			videosToInsert = nil
		}

		if apiResponse.NextPageToken == "" {
			break
		}

		nextPageToken = apiResponse.NextPageToken

		time.Sleep(time.Second * 10)
	}
}
