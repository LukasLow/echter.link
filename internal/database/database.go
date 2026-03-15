package database

import (
	"database/sql"
	"log"
	"os"

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
	`

	_, err = DB.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ V1.1: Go-Backend & SQLite Datenbank initialisiert")
}
