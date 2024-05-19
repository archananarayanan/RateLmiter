package main

import (
	_ "RateLmiter/docs"
	r "RateLmiter/routes"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title RateLimiter APIs
// @version 1.0
// @description Testing Rate Limiter.
// @termsOfService http://swagger.io/terms/
// @BasePath /api/v1

// @schemes http
func main() {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/api/v1/requestLimit", r.RequestLimit)
	router.GET("/api/v1/criticalRequestLimit", r.CriticalRequestLimit)
	router.Run("localhost:8080")
}