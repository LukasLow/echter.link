package handlers

import (
	"database/sql"
	"fmt"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Anfrage. Bitte überprüfe deine Eingaben."})
		return
	}

	// Smart URL Auto-Complete: Auto-https:// wenn .com/.de/.org etc.
	normalizedURL := NormalizeURL(request.OriginalURL)

	// V1.10: Validierung von URLs (Anti-Phishing Check)
	if !IsValidURL(normalizedURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige URL. Die URL muss mit http:// oder https:// beginnen."})
		return
	}

	shortURL := models.ShortURL{
		ID:          uuid.New().String(),
		OriginalURL: normalizedURL,
		CreatedAt:   time.Now(),
	}

	if request.CustomCode != "" {
		// Detaillierte Validierung mit spezifischer Fehlermeldung
		if len(request.CustomCode) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Der eigene Kurzlink-Code ist zu kurz. Mindestens 3 Zeichen erforderlich. Du hast nur " + fmt.Sprintf("%d", len(request.CustomCode)) + " Zeichen eingegeben."})
			return
		}
		if len(request.CustomCode) > 32 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Der eigene Kurzlink-Code ist zu lang. Maximal 32 Zeichen erlaubt. Du hast " + fmt.Sprintf("%d", len(request.CustomCode)) + " Zeichen eingegeben."})
			return
		}
		for _, r := range request.CustomCode {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültiges Zeichen im Kurzlink-Code: '" + string(r) + "'. Erlaubt sind nur: Buchstaben (a-z, A-Z), Zahlen (0-9), Bindestrich (-) und Unterstrich (_)."})
				return
			}
		}
		var existingID string
		err := database.DB.QueryRow("SELECT id FROM short_urls WHERE short_code = ?", request.CustomCode).Scan(&existingID)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Der Kurzlink-Code '" + request.CustomCode + "' ist bereits vergeben. Bitte wähle einen anderen."})
			return
		}
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Prüfen der Verfügbarkeit. Bitte versuche es später erneut."})
			return
		}
		shortURL.ShortCode = request.CustomCode
	} else {
		shortCode, err := GenerateUniqueShortCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen des Kurzlinks. Bitte versuche es später erneut."})
			return
		}
		shortURL.ShortCode = shortCode
	}

	if request.ExpiresIn > 0 {
		// Anonymous users limited to 7 days (168 hours)
		if request.ExpiresIn > 168 {
			request.ExpiresIn = 168
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Speichern des Kurzlinks. Bitte versuche es später erneut."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url":  database.GetDomain() + "/" + shortURL.ShortCode,
		"short_code": shortURL.ShortCode,
		"expires_at": shortURL.ExpiresAt,
	})
}

func GenerateShortCode() string {
	return strings.ReplaceAll(uuid.New().String()[:8], "-", "")
}

func GenerateUniqueShortCode() (string, error) {
	for i := 0; i < 10; i++ {
		code := GenerateShortCode()
		var existingID string
		err := database.DB.QueryRow("SELECT id FROM short_urls WHERE short_code = ?", code).Scan(&existingID)
		if err == sql.ErrNoRows {
			return code, nil
		}
		if err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("failed to generate unique code after 10 attempts")
}

func IsValidCustomCode(code string) bool {
	if len(code) < 3 || len(code) > 32 {
		return false
	}
	for _, r := range code {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

func RedirectShortCode(c *gin.Context) {
	shortCode := c.Param("code")

	var originalURL string
	var expiresAt *time.Time
	err := database.DB.QueryRow("SELECT original_url, expires_at FROM short_urls WHERE short_code = ?", shortCode).Scan(&originalURL, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, "Kurzlink nicht gefunden. Der Link existiert nicht oder wurde gelöscht.")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen des Kurzlinks. Bitte versuche es später erneut."})
		return
	}

	if expiresAt != nil && time.Now().After(*expiresAt) {
		c.String(http.StatusGone, "Dieser Kurzlink ist abgelaufen und nicht mehr verfügbar.")
		return
	}

	_, err = database.DB.Exec("UPDATE short_urls SET clicks = clicks + 1 WHERE short_code = ?", shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Aktualisieren der Klick-Statistik."})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
