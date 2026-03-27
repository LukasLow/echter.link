# echter.link

Ein moderner, selbst hostbarer URL-Shortener mit deutscher Benutzeroberfläche, Dark Mode und einfacher Bedienung. Gebaut mit Go, Gin und SQLite.

![Version](https://img.shields.io/badge/version-1.3.1-blue)
![Go](https://img.shields.io/badge/go-1.21-green)
![Docker](https://img.shields.io/badge/docker-ready-blue)

## ✨ Features

- **🔗 Kurzlinks erstellen:** Wähle zwischen zufälligen oder eigenen Codes
- **⏰ Ablaufdatum:** Links automatisch nach 1h, 24h, 3d oder 7d ablaufen lassen
- **🌙 Dark Mode:** Automatisch, Hell oder Dunkel - mit System-Erkennung
- **🇩🇪 Deutsche UI:** Vollständig übersetzt und benutzerfreundlich
- **📱 Responsive:** Funktioniert auf Desktop, Tablet und Smartphone
- **🐳 Docker:** Einfache Installation mit Docker Compose
- **💾 SQLite:** Keine Datenbank-Setup nötig, alles in einer Datei
- **🚀 Schnell:** Optimiert für kleine Server (0.5 vCPU, 500MB RAM)

## 🚀 Schnellstart

### Mit Docker (Empfohlen)

```bash
# Repository klonen
git clone https://github.com/LukasLow/echter.link.git
cd echter.link

# Server starten
docker compose up -d
```

### Mit Podman (Rootless, sicherer)

```bash
# Repository klonen
git clone https://github.com/LukasLow/echter.link.git
cd echter.link

# Mit podman-compose
podman-compose up -d

# Oder Podman native
podman build -t echter-link .
podman run -d -p 8080:8080 -v ./data:/root/data:z echter-link
```

Die Webseite ist dann unter `http://localhost:8080` erreichbar.

### Lokal entwickeln

```bash
# Go 1.21+ benötigt
git clone https://github.com/LukasLow/echter.link.git
cd echter.link
go mod tidy
go run cmd/server.go
```

## � Projektstruktur

```
echter.link/
├── cmd/
│   └── server.go              # Hauptanwendung
├── internal/
│   ├── database/
│   │   └── database.go        # Datenbank-Setup
│   ├── handlers/
│   │   ├── url.go            # API-Logik
│   │   └── web.go            # Web-Interface
│   └── models/
│       └── shorturl.go        # Datenmodelle
├── Dockerfile                 # Docker-Build
├── docker-compose.yml         # Docker-Compose Konfig
├── go.mod                     # Go-Abhängigkeiten
├── CHANGELOG.md              # Versionshistorie
└── README.md                 # Diese Datei
```

## 🔧 Konfiguration

Über Umgebungsvariablen:

| Variable | Standard | Beschreibung |
|----------|----------|--------------|
| `GIN_MODE` | `release` | Gin Framework Modus |
| `DB_PATH` | `./echter.link.sqlite` | SQLite Datenbankpfad |
| `DOMAIN` | `http://localhost:8080` | Basis-URL für Kurzlinks |

## 📊 API

### Kurzlink erstellen

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "original_url": "https://example.com",
    "custom_code": "mein-link",
    "expires_in": 24
  }'
```

**Antwort:**
```json
{
  "short_url": "http://localhost:8080/mein-link",
  "short_code": "mein-link",
  "expires_at": "2026-03-28T20:00:00Z"
}
```

### Endpunkte

| Methode | URL | Beschreibung |
|---------|-----|--------------|
| `GET` | `/` | Web-Interface |
| `POST` | `/api/shorten` | Kurzlink erstellen |
| `GET` | `/:code` | Weiterleitung zum Original |

## 🗄️ Datenbank

SQLite mit einer Tabelle:

```sql
CREATE TABLE short_urls (
    id TEXT PRIMARY KEY,
    short_code TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    clicks INTEGER DEFAULT 0,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🔒 Sicherheit

- **URL-Validierung:** Nur `http://` und `https://` erlaubt
- **Eingabe-Validierung:** Custom Codes: 3-32 Zeichen, alphanumerisch + `-` und `_`
- **Zeitlimit:** Anonyme Nutzer maximal 7 Tage
- **Graceful Shutdown:** Saubere Beendigung ohne Datenverlust

## 📈 Performance

- **Speicher:** ~250MB Basis, 500MB Limit
- **CPU:** Optimiert für 0.5 vCPU
- **Antwortzeit:** <100ms für API-Aufrufe
- **Container-Größe:** ~20MB komprimiert

## � Versionen

Siehe [CHANGELOG.md](CHANGELOG.md) für die vollständige Versionshistorie.

## 📄 Lizenz

Open Source - siehe LICENSE Datei.

---

**Gebaut mit ❤️ mit Go, Gin & SQLite**
