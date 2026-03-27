package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Environment Setup: Domain als ENV Host + localhost:8080
func GetDomain() string {
	if domain := os.Getenv("DOMAIN"); domain != "" {
		return domain
	}
	return "http://localhost:8080"
}

func InitDB() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./echter.link.sqlite"
	}

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Configure connection pool for concurrent access
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	createTables := `
	CREATE TABLE IF NOT EXISTS short_urls (
		id TEXT PRIMARY KEY,
		short_code TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		expires_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_short_code ON short_urls(short_code);

	CREATE TABLE IF NOT EXISTS admin_users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = DB.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}

	// Create default admin user if not exists
	createDefaultAdmin()

	log.Println("✅ V1.1: Go-Backend & SQLite Datenbank initialisiert")
}

// Create default admin from environment variables
func createDefaultAdmin() {
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPass == "" {
		adminPass = "admin123"
	}

	// Check if admin already exists
	var existingID string
	err := DB.QueryRow("SELECT id FROM admin_users WHERE username = ?", adminUser).Scan(&existingID)
	if err == nil {
		log.Println("✅ Admin-Benutzer existiert bereits")
		return
	}
	if err != sql.ErrNoRows {
		log.Println("⚠️ Fehler beim Prüfen des Admin-Benutzers:", err)
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	if err != nil {
		log.Println("⚠️ Fehler beim Hashen des Admin-Passworts:", err)
		return
	}

	// Insert admin
	_, err = DB.Exec("INSERT INTO admin_users (id, username, password_hash) VALUES (?, ?, ?)",
		uuid.New().String(), adminUser, string(hash))
	if err != nil {
		log.Println("⚠️ Fehler beim Erstellen des Admin-Benutzers:", err)
		return
	}

	log.Printf("✅ Admin-Benutzer '%s' erstellt (Passwort aus ADMIN_PASSWORD oder Standard: admin123)", adminUser)
}
