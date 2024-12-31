package router

import (
	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/Fampay-Assignment/internal/videos"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	authGroup := router.Group("/fam")
	videos.RegisterRoutes(authGroup)

	return router
}
