package videos

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/Fampay-Assignment/internal/db/postgres_db"
	"github.com/harsh082ip/Fampay-Assignment/internal/helpers"
	"github.com/harsh082ip/Fampay-Assignment/internal/models"
)

func GetVideos(c *gin.Context) {
	limit := c.Query("limit")
	pageToken := c.Query("pageToken")

	// Define structure for the response
	var videos []models.Video
	var pageTokenRes int

	// Validate limit parameter
	if limit == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "limit is required",
		})
		return
	}

	// Convert limit to integer
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "Invalid limit value",
		})
		return
	}

	if pageToken == "" {
		// First query: Get initial set of videos
		videoQuery := fmt.Sprintf(`
            SELECT *       
            FROM (
                SELECT * 
                FROM videos 
                ORDER BY id
                LIMIT %v
            ) AS limited_videos 
            ORDER BY published_at DESC;
        `, limitInt)

		result := postgres_db.DB.Raw(videoQuery).Scan(&videos)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
				"error":  "Cannot fetch videos from DB, " + result.Error.Error(),
			})
			return
		}

		// Get max ID from the results for next page token
		if len(videos) > 0 {

			pageTokenQuery := fmt.Sprintf(`
				SELECT MAX(id) AS max_id                 
                FROM (
                        SELECT id 
                        FROM videos 
                        LIMIT %v
                ) AS limited_ids;
			`, limitInt)

			// maxIDQuery := `SELECT MAX(id) FROM videos`
			if err := postgres_db.DB.Raw(pageTokenQuery).Scan(&pageTokenRes); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "Internal Server Error",
					"error":  "Cannot fetch max ID from DB, " + err.Error.Error(),
				})
				return
			}
		}

		base64PageToken, err := helpers.IntToBase64(pageTokenRes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
				"error":  "error encoding to base64: " + err.Error(),
			})
			return
		}

		apiResponse := models.ApiResponse{
			PageToken: base64PageToken,
			Videos:    videos,
		}

		log.Println(len(apiResponse.Videos))
		c.JSON(http.StatusOK, apiResponse)
		return
	}

	// Handle pagination with pageToken
	pageTokenInt, err := helpers.Base64ToInt(pageToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
			"error":  "Invalid page token: " + err.Error(),
		})
		return
	}
	log.Println("int:", pageTokenInt)

	// Get next page of videos
	videosQuery := fmt.Sprintf(`
        SELECT *
        FROM (
            SELECT *
            FROM videos 
            WHERE id > %v
            ORDER BY id
            LIMIT %v
        ) subq
        ORDER BY published_at DESC;
    `, pageTokenInt, limitInt)

	result := postgres_db.DB.Raw(videosQuery).Scan(&videos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
			"error":  "Cannot fetch videos from DB, " + result.Error.Error(),
		})
		return
	}

	// Get next page token
	if len(videos) > 0 {
		nextPageQuery := fmt.Sprintf(`
            SELECT MAX(id)
            FROM (
                SELECT id
                FROM videos
                WHERE id > %v
                ORDER BY id
                LIMIT %v
            ) AS ids;
        `, pageTokenInt, limitInt)

		if err := postgres_db.DB.Raw(nextPageQuery).Scan(&pageTokenRes); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Server Error",
				"error":  "Cannot fetch next page token from DB, " + err.Error.Error(),
			})
			return
		}
	}

	base64PageToken, err := helpers.IntToBase64(pageTokenRes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
			"error":  "error encoding to base64: " + err.Error(),
		})
		return
	}

	apiResponse := models.ApiResponse{
		PageToken: base64PageToken,
		Videos:    videos,
	}

	log.Println(len(apiResponse.Videos))
	c.JSON(http.StatusOK, apiResponse)
}
