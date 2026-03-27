# Changelog

Alle wichtigen Änderungen an diesem Projekt werden hier dokumentiert.

Das Format folgt [Keep a Changelog](https://keepachangelog.com/de/1.0.0/) und [Semantic Versioning](https://semver.org/lang/de/).

---

## [1.3.1] - 2026-03-27

### Verbessert (Changed)
- **Deutsche Übersetzung:** Alle Fehlermeldungen und UI-Texte ins Deutsche übersetzt
- **Benutzerfreundlichere Texte:** Klare, verständliche Formulierungen für alle Fehlermeldungen
- **Detaillierte Validierung:** Fehlermeldungen zeigen genau, welches Zeichen ungültig ist oder wie viele Zeichen eingegeben wurden
- **UI-Labels optimiert:** "Eigener Code (optional)" statt technischem Fachjargon

### Behoben (Fixed)
- **Kurzlink-Weiterleitung:** Entfernt '#' aus der URL-Generierung, die die Weiterleitung blockiert hat

---

## [1.3.0] - 2026-03-27

### Verbessert (Changed)
- **Ablaufdatum-Auswahl:** Einfaches Dropdown statt Zahlen-Eingabe (1h, 24h, 3d, 7d)
- **Zeitlimit für anonyme Nutzer:** Maximale Gültigkeit auf 7 Tage (168 Stunden) begrenzt

---

## [1.2.0] - 2026-03-27

### Neu (Added)
- **Dark Mode:** Automatische System-Erkennung + manuelle Umschaltung
- **Theme-Toggle:** Wählbar zwischen Hell ☀️, Dunkel 🌙 und System 🖥️
- **Theme-Speicherung:** Auswahl wird im Browser gespeichert
- **CSS Variablen:** Nahtlose Übergänge zwischen Themes

---

## [1.1.0] - 2026-03-27

### Neu (Added)
- **Modernes UI-Design:** Gradient-Hintergrund mit weißer Card
- **Footer mit Version:** Versions-Badge und strukturierte Links
- **Changelog-Link:** Direkter Link zur Versionshistorie
- **Enter-Taste Support:** Schnelles Erstellen mit Enter

### Behoben (Fixed)
- **Race Condition:** Ablaufprüfung jetzt vor Klick-Zählung
- **URL-Validierung:** Reihenfolge korrigiert (Normalisierung → Validierung)
- **Custom Code Prüfung:** Datenbankfehler werden korrekt behandelt
- **Code-Kollisionen:** Retry-Logik für eindeutige Codes
- **Eingabevalidierung:** Zeichenlimits (3-32) und erlaubte Zeichen enforced

### Verbessert (Changed)
- **Redesign:** Moderner Look mit Gradient und Card-Layout
- **Footer überarbeitet:** Versionsanzeige und saubere Link-Struktur
- **Performance:** Snow/Mouse-Glow Effekte entfernt für bessere Performance

### Sicherheit (Security)
- **Graceful Shutdown:** Saubere Server-Beendigung
- **Connection Pool:** 25 Datenbank-Verbindungen konfiguriert

---

## [1.0.0] - 2026-03-26

### Erste Veröffentlichung (Initial Release)
- **URL-Shortening:** Zufällige und benutzerdefinierte Kurzlinks
- **Ablaufdaten:** Links mit Zeitlimit (1-720 Stunden)
- **Click-Tracking:** Klickzähler für jeden Link
- **SQLite Datenbank:** Dateibasierte Persistenz
- **Docker Support:** Multi-Stage Build für kleine Images
- **REST API:** `/api/shorten` Endpunkt
- **Web UI:** Eingebettetes HTML/CSS/JS Interface

---
