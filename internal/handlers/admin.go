package handlers

import (
	"net/http"
	"os"
	"time"

	"echter.link/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "default-secret-change-in-production"))

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// AdminLogin handles admin authentication
func AdminLogin(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Benutzername und Passwort erforderlich"})
		return
	}

	// Verify credentials
	var passwordHash string
	err := database.DB.QueryRow("SELECT password_hash FROM admin_users WHERE username = ?", request.Username).Scan(&passwordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültige Anmeldedaten"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültige Anmeldedaten"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": request.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen der Sitzung"})
		return
	}

	// Set HTTP-only cookie
	c.SetCookie("admin_session", tokenString, 86400, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Erfolgreich angemeldet"})
}

// AdminLogout handles logout
func AdminLogout(c *gin.Context) {
	c.SetCookie("admin_session", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Erfolgreich abgemeldet"})
}

// AdminStats returns dashboard statistics
func AdminStats(c *gin.Context) {
	var totalLinks, totalClicks, activeLinks int

	// Total links
	database.DB.QueryRow("SELECT COUNT(*) FROM short_urls").Scan(&totalLinks)

	// Total clicks
	database.DB.QueryRow("SELECT COALESCE(SUM(clicks), 0) FROM short_urls").Scan(&totalClicks)

	// Active links (not expired)
	database.DB.QueryRow("SELECT COUNT(*) FROM short_urls WHERE expires_at IS NULL OR expires_at > ?", time.Now()).Scan(&activeLinks)

	c.JSON(http.StatusOK, gin.H{
		"total_links":  totalLinks,
		"total_clicks": totalClicks,
		"active_links": activeLinks,
	})
}

// AdminLinks returns all links with pagination
func AdminLinks(c *gin.Context) {
	page := 1
	limit := 50
	if p := c.Query("page"); p != "" {
		// Simple pagination logic would go here
		_ = p
	}

	rows, err := database.DB.Query(`
		SELECT short_code, original_url, clicks, expires_at, created_at 
		FROM short_urls 
		ORDER BY created_at DESC 
		LIMIT ?
	`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Laden der Links"})
		return
	}
	defer rows.Close()

	var links []gin.H
	for rows.Next() {
		var link struct {
			ShortCode   string     `json:"short_code"`
			OriginalURL string     `json:"original_url"`
			Clicks      int        `json:"clicks"`
			ExpiresAt   *time.Time `json:"expires_at"`
			CreatedAt   time.Time  `json:"created_at"`
		}
		rows.Scan(&link.ShortCode, &link.OriginalURL, &link.Clicks, &link.ExpiresAt, &link.CreatedAt)
		links = append(links, gin.H{
			"short_code":   link.ShortCode,
			"original_url": link.OriginalURL,
			"clicks":       link.Clicks,
			"expires_at":   link.ExpiresAt,
			"created_at":   link.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"links": links,
		"page":  page,
		"limit": limit,
	})
}

// AdminDeleteLink deletes a short link
func AdminDeleteLink(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kein Kurzlink-Code angegeben"})
		return
	}

	_, err := database.DB.Exec("DELETE FROM short_urls WHERE short_code = ?", code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Löschen des Links"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link erfolgreich gelöscht"})
}

// AuthRequired middleware validates JWT token
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("admin_session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Nicht autorisiert"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültige Sitzung"})
			c.Abort()
			return
		}

		c.Next()
	}
}
