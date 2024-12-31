package videos

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/videos", GetVideos)
	router.GET("/videos/:id", GetVideoByID)
}
