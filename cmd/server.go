package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"echter.link/internal/database"
	"echter.link/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/api/shorten", handlers.CreateShortURL)
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Admin routes
	r.GET("/admin", handlers.AdminLoginHandler)
	r.POST("/admin/login", handlers.AdminLogin)
	r.GET("/admin/logout", handlers.AdminLogout)

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(handlers.AuthRequired())
	{
		admin.GET("/dashboard", handlers.AdminDashboardHandler)
		admin.GET("/api/stats", handlers.AdminStats)
		admin.GET("/api/links", handlers.AdminLinks)
		admin.DELETE("/api/links/:code", handlers.AdminDeleteLink)
	}

	// Handle redirects for short URLs
	r.GET("/:code", handlers.RedirectShortCode)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Graceful shutdown handling
	go func() {
		log.Println("🚀 echter.link V1.1 starting on :8080")
		log.Println("✅ Features: Go-Backend, SQLite, Random Codes, Custom Codes, URL Validation, Expiration")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
