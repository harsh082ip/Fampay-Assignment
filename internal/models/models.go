// models/models.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type YoutubeApiResponse struct {
	NextPageToken string      `json:"nextPageToken"`
	Items         []VideoItem `json:"items"`
}

type VideoItem struct {
	ID      VideoID      `json:"id"`
	Snippet VideoSnippet `json:"snippet"`
}

type VideoID struct {
	VideoID string `json:"videoId"`
}

type VideoSnippet struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Thumbnails  VideoThumbnails `json:"thumbnails"`
	PublishTime string          `json:"publishTime"`
}

type VideoThumbnails struct {
	Default DefaultRes `json:"default"`
	Medium  MediumRes  `json:"medium"`
	High    HighRes    `json:"high"`
}

type DefaultRes struct {
	URL string `json:"url"`
}

type MediumRes struct {
	URL string `json:"url"`
}

type HighRes struct {
	URL string `json:"url"`
}

type Video struct {
	gorm.Model
	VideoID       string    `gorm:"not null;primaryKey" json:"videoId"`
	Title         string    `gorm:"not null" json:"title"`
	Description   *string   `gorm:"type:text" json:"description"`
	PublishedAt   time.Time `gorm:"type:timestamp;not null" json:"publishedAt"`
	ThumbnailURLs *string   `gorm:"type:text" json:"thumbnails"`
}
