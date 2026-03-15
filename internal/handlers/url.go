package handlers

import (
	"net/http"
	"strings"
	"time"

	"echter.link/internal/database"
	"echter.link/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// V1.10: Validierung von URLs (Anti-Phishing Check)
func IsValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// Smart URL Auto-Complete: Auto-https:// wenn .com/.de/.org etc.
func NormalizeURL(url string) string {
	// Remove leading/trailing whitespace
	url = strings.TrimSpace(url)

	// If URL doesn't have protocol, add https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}

	return url
}

func CreateShortURL(c *gin.Context) {
	var request struct {
		OriginalURL string `json:"original_url" binding:"required"`
		CustomCode  string `json:"custom_code"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// V1.10: Validierung von URLs (Anti-Phishing Check)
	if !IsValidURL(request.OriginalURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format - must start with http:// or https://"})
		return
	}

	// Smart URL Auto-Complete: Auto-https:// wenn .com/.de/.org etc.
	normalizedURL := NormalizeURL(request.OriginalURL)

	shortURL := models.ShortURL{
		ID:          uuid.New().String(),
		OriginalURL: normalizedURL,
		CreatedAt:   time.Now(),
	}

	if request.CustomCode != "" {
		var existingID string
		err := database.DB.QueryRow("SELECT id FROM short_urls WHERE short_code = ?", request.CustomCode).Scan(&existingID)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Custom code already taken"})
			return
		}
		shortURL.ShortCode = request.CustomCode
	} else {
		shortURL.ShortCode = GenerateShortCode()
	}

	if request.ExpiresIn > 0 {
		if request.ExpiresIn > 720 {
			request.ExpiresIn = 720
		}
		expiresAt := time.Now().Add(time.Duration(request.ExpiresIn) * time.Hour)
		shortURL.ExpiresAt = &expiresAt
	}

	var err error
	if shortURL.ExpiresAt != nil {
		_, err = database.DB.Exec("INSERT INTO short_urls (id, short_code, original_url, expires_at) VALUES (?, ?, ?, ?)",
			shortURL.ID, shortURL.ShortCode, shortURL.OriginalURL, shortURL.ExpiresAt)
	} else {
		_, err = database.DB.Exec("INSERT INTO short_urls (id, short_code, original_url) VALUES (?, ?, ?)",
			shortURL.ID, shortURL.ShortCode, shortURL.OriginalURL)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url":  database.GetDomain() + "/#" + shortURL.ShortCode,
		"short_code": shortURL.ShortCode,
		"expires_at": shortURL.ExpiresAt,
	})
}

func GenerateShortCode() string {
	return strings.ReplaceAll(uuid.New().String()[:8], "-", "")
}

func RedirectShortCode(c *gin.Context) {
	shortCode := c.Param("code")

	result, err := database.DB.Exec("UPDATE short_urls SET clicks = clicks + 1 WHERE short_code = ?", shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update clicks"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.String(http.StatusNotFound, "Short URL not found")
		return
	}

	var originalURL string
	var expiresAt *time.Time
	err = database.DB.QueryRow("SELECT original_url, expires_at FROM short_urls WHERE short_code = ?", shortCode).Scan(&originalURL, &expiresAt)
	if err != nil {
		c.String(http.StatusNotFound, "Short URL not found")
		return
	}

	if expiresAt != nil && time.Now().After(*expiresAt) {
		c.String(http.StatusGone, "Short URL has expired")
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
