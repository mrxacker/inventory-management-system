package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/config"
)

func SetupRouter(productHandler *ProductHandler, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "product-service",
			"time":    time.Now(),
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// API routes
	// v1 := router.Group("/api/v1")
	// {
	// 	products := v1.Group("/products")
	// 	{
	// 		products.GET("", productHandler.GetProducts)
	// 		products.GET("/:id", productHandler.GetProduct)
	// 		products.POST("", productHandler.CreateProduct)
	// 		products.PUT("/:id", productHandler.UpdateProduct)
	// 		products.DELETE("/:id", productHandler.DeleteProduct)
	// 	}
	// }

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
