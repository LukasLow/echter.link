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

## 🖥️ Server-Setup (Debian)

Komplette Anleitung für einen frischen Debian-Server (1 CPU, 1GB RAM, 10GB Storage). 
**Kein git nötig** - das Image kommt fertig gebaut von GitHub Container Registry.

### 1. System aktualisieren

```bash
# Als root oder mit sudo
apt update && apt upgrade -y
apt install -y curl
```

### 2. Podman installieren

```bash
apt install -y podman

# Überprüfen
podman --version
```

### 3. Datenverzeichnis anlegen

```bash
# Persistente Daten für SQLite
mkdir -p /opt/echter.link/data
chmod 755 /opt/echter.link/data
```

### 4. Container starten (fertiges Image)

```bash
# Pre-built Image von GitHub Container Registry laden
podman run -d \
  --name echter-link \
  -p 8080:8080 \
  -v /opt/echter.link/data:/root/data:z \
  -e GIN_MODE=release \
  -e DB_PATH=/root/data/echter.link.sqlite \
  -e DOMAIN=https://echter.link \
  --restart unless-stopped \
  ghcr.io/lukaslow/echter.link:latest
```


### 6. SSL/HTTPS mit Caddy (leichtgewichtig)

Caddy ist perfekt für kleine Server - automatisches HTTPS, sehr wenig RAM/CPU.

```bash
# Caddy installieren
apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt install -y caddy

# Caddy konfigurieren
cat > /etc/caddy/Caddyfile << 'EOF'
echter.link {
    reverse_proxy localhost:8080
}
EOF

# Caddy starten
systemctl enable caddy
systemctl start caddy
```

### 7. Automatische Updates mit Watchtower

```bash
# Watchtower Container starten (prüft alle 24h auf Updates)
podman run -d \
  --name watchtower \
  -v /run/user/$(id - u)/podman/podman.sock:/var/run/docker.sock:z \
  --restart unless-stopped \
  containrrr/watchtower \
  --cleanup \
  --interval 86400 \
  echter-link
```

### 8. Als Systemd-Service (optional)

Für automatischen Start beim Booten:

```bash
# Podman-Container als Systemd-Service generieren
podman generate systemd --name echter-link --files

# Service-Datei verschieben
mv container-echter-link.service /etc/systemd/system/

# Aktivieren und starten
systemctl daemon-reload
systemctl enable container-echter-link.service
systemctl start container-echter-link.service
```

### 9. Backup der Datenbank

```bash
# Cron-Job für tägliches Backup
cat > /etc/cron.daily/backup-echter-link << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/echter.link/backups"
mkdir -p "$BACKUP_DIR"
cp "/opt/echter.link/data/echter.link.sqlite" "$BACKUP_DIR/echter.link-$(date +%Y%m%d).sqlite"
# Alte Backups löschen (älter als 7 Tage)
find "$BACKUP_DIR" -name "echter.link-*.sqlite" -mtime +7 -delete
EOF

chmod +x /etc/cron.daily/backup-echter-link
```

### Fertig!

Dein echter.link ist jetzt unter `https://echter.link` erreichbar mit:
- ✅ Automatischem SSL (Caddy)
- ✅ Persistenten Daten (SQLite)
- ✅ Automatischen Updates (Watchtower)
- ✅ Täglichen Backups

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
