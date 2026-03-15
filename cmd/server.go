package main

import (
	"echter.link/internal/database"
	"echter.link/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/api/shorten", handlers.CreateShortURL)
	r.GET("/", handlers.HomeHandler)
	
	// Handle redirects for short URLs
	r.GET("/:code", handlers.RedirectShortCode)
	
	// Handle hash-based redirects for frontend
	r.GET("/#/:code", handlers.RedirectShortCode)

	log.Println("🚀 echter.link V1.1 starting on :8080")
	log.Println("✅ Features: Go-Backend, SQLite, Random Codes, Custom Codes, URL Validation, Expiration")
	r.Run(":8080")
}
