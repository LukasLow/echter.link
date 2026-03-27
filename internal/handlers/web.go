package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>echter.link - URL Shortener</title>
    <style>
        :root {
            --bg-gradient-start: #667eea;
            --bg-gradient-end: #764ba2;
            --card-bg: white;
            --card-shadow: 0 20px 60px rgba(0,0,0,0.3);
            --text-primary: #333;
            --text-secondary: #666;
            --input-bg: #f8f9fa;
            --input-border: #e9ecef;
            --input-focus: #667eea;
            --result-bg: #f8f9fa;
            --result-border: #667eea;
            --btn-primary-bg: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            --btn-copy-border: #667eea;
            --btn-copy-color: #667eea;
            --btn-copy-hover-bg: #667eea;
            --btn-copy-hover-color: white;
            --error-bg: #fee;
            --error-border: #e74c3c;
            --error-color: #c33;
            --link-color: #667eea;
            --footer-bg: rgba(0,0,0,0.2);
            --footer-text: white;
            --toggle-bg: rgba(255,255,255,0.2);
            --toggle-active: white;
        }

        [data-theme="dark"] {
            --bg-gradient-start: #1a1a2e;
            --bg-gradient-end: #16213e;
            --card-bg: #1e1e2e;
            --card-shadow: 0 20px 60px rgba(0,0,0,0.5);
            --text-primary: #eaeaea;
            --text-secondary: #a0a0a0;
            --input-bg: #2a2a3a;
            --input-border: #3a3a4a;
            --input-focus: #8b5cf6;
            --result-bg: #2a2a3a;
            --result-border: #8b5cf6;
            --btn-primary-bg: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
            --btn-copy-border: #8b5cf6;
            --btn-copy-color: #8b5cf6;
            --btn-copy-hover-bg: #8b5cf6;
            --btn-copy-hover-color: white;
            --error-bg: #3a1f1f;
            --error-border: #ef4444;
            --error-color: #ef4444;
            --link-color: #8b5cf6;
            --footer-bg: rgba(0,0,0,0.4);
            --footer-text: #d0d0d0;
            --toggle-bg: rgba(255,255,255,0.1);
            --toggle-active: #8b5cf6;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            background: linear-gradient(135deg, var(--bg-gradient-start) 0%, var(--bg-gradient-end) 100%);
            color: var(--text-primary);
            line-height: 1.6;
            transition: background 0.3s ease;
        }

        .main-container {
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }

        .card {
            background: var(--card-bg);
            border-radius: 20px;
            box-shadow: var(--card-shadow);
            padding: 50px 40px;
            width: 100%;
            max-width: 520px;
            text-align: center;
            transition: background 0.3s ease, box-shadow 0.3s ease;
        }

        .card-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 20px;
        }

        .theme-toggle {
            display: flex;
            background: var(--toggle-bg);
            border-radius: 25px;
            padding: 4px;
            gap: 4px;
        }

        .theme-btn {
            border: none;
            background: transparent;
            color: var(--footer-text);
            padding: 6px 12px;
            border-radius: 20px;
            cursor: pointer;
            font-size: 0.85rem;
            transition: all 0.2s;
            display: flex;
            align-items: center;
            gap: 4px;
        }

        .theme-btn:hover {
            background: rgba(255,255,255,0.1);
        }

        .theme-btn.active {
            background: var(--toggle-active);
            color: var(--text-primary);
        }

        [data-theme="dark"] .theme-btn.active {
            color: white;
        }

        .logo {
            font-size: 2.5rem;
            font-weight: 800;
            background: var(--btn-primary-bg);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 8px;
        }

        .tagline {
            color: var(--text-secondary);
            font-size: 0.95rem;
            margin-bottom: 35px;
        }

        .input-group {
            display: flex;
            margin-bottom: 20px;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }

        .input-prefix {
            background: var(--input-bg);
            padding: 16px 14px;
            font-size: 0.9rem;
            color: var(--text-secondary);
            border: none;
            font-weight: 500;
            transition: background 0.3s ease;
        }

        input[type="url"] {
            flex: 1;
            padding: 16px;
            border: none;
            font-size: 1rem;
            outline: none;
            background: var(--card-bg);
            color: var(--text-primary);
            transition: background 0.3s ease, color 0.3s ease;
        }

        .options {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
            margin-bottom: 25px;
        }

        .option-group {
            text-align: left;
        }

        .option-group label {
            display: block;
            font-size: 0.8rem;
            color: var(--text-secondary);
            margin-bottom: 6px;
            font-weight: 500;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .option-group input {
            width: 100%;
            padding: 12px 14px;
            border: 2px solid var(--input-border);
            border-radius: 10px;
            font-size: 0.95rem;
            outline: none;
            background: var(--card-bg);
            color: var(--text-primary);
            transition: border-color 0.2s, background 0.3s ease, color 0.3s ease;
        }

        .option-group input:focus {
            border-color: var(--input-focus);
        }

        .btn-primary {
            width: 100%;
            padding: 16px 28px;
            background: var(--btn-primary-bg);
            color: white;
            border: none;
            border-radius: 12px;
            font-size: 1.05rem;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 25px rgba(102, 126, 234, 0.5);
        }

        .result {
            margin-top: 25px;
            padding: 20px;
            background: var(--result-bg);
            border-radius: 12px;
            border-left: 4px solid var(--result-border);
            text-align: left;
            transition: background 0.3s ease;
        }

        .result a {
            color: var(--link-color);
            word-break: break-all;
            text-decoration: none;
            font-weight: 500;
        }

        .result a:hover {
            text-decoration: underline;
        }

        .btn-copy {
            margin-top: 12px;
            padding: 10px 20px;
            background: var(--card-bg);
            border: 2px solid var(--btn-copy-border);
            color: var(--btn-copy-color);
            border-radius: 8px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.2s;
        }

        .btn-copy:hover {
            background: var(--btn-copy-hover-bg);
            color: var(--btn-copy-hover-color);
        }

        /* Footer */
        footer {
            background: var(--footer-bg);
            backdrop-filter: blur(10px);
            color: var(--footer-text);
            padding: 25px 20px;
            text-align: center;
            transition: background 0.3s ease, color 0.3s ease;
        }

        .footer-content {
            max-width: 800px;
            margin: 0 auto;
        }

        .version-badge {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            background: rgba(255,255,255,0.15);
            padding: 6px 14px;
            border-radius: 20px;
            font-size: 0.85rem;
            margin-bottom: 12px;
        }

        .version-badge .label {
            opacity: 0.8;
        }

        .version-badge .version {
            font-weight: 600;
            background: rgba(255,255,255,0.2);
            padding: 2px 8px;
            border-radius: 4px;
        }

        .footer-links {
            margin: 15px 0;
            font-size: 0.9rem;
        }

        .footer-links a {
            color: rgba(255,255,255,0.9);
            text-decoration: none;
            margin: 0 12px;
            transition: color 0.2s;
        }

        .footer-links a:hover {
            color: white;
            text-decoration: underline;
        }

        [data-theme="dark"] .footer-links a {
            color: rgba(208,208,208,0.9);
        }

        [data-theme="dark"] .footer-links a:hover {
            color: white;
        }

        .footer-tech {
            font-size: 0.8rem;
            opacity: 0.7;
            margin-top: 10px;
        }

        .error {
            background: var(--error-bg);
            border-left-color: var(--error-border);
            color: var(--error-color);
        }

        @media (max-width: 480px) {
            .card {
                padding: 35px 25px;
            }
            .card-header {
                flex-direction: column;
                align-items: center;
                gap: 15px;
            }
            .theme-toggle {
                order: -1;
            }
            .logo {
                font-size: 2rem;
            }
            .options {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="main-container">
        <div class="card">
            <div class="card-header">
                <div></div>
                <div class="theme-toggle">
                    <button class="theme-btn" data-theme="light" title="Hell">☀️</button>
                    <button class="theme-btn" data-theme="dark" title="Dunkel">🌙</button>
                    <button class="theme-btn active" data-theme="auto" title="System">🖥️</button>
                </div>
            </div>

            <h1 class="logo">echter.link</h1>
            <p class="tagline">Kurze Links. Echte Einfachheit.</p>
            
            <div class="input-group">
                <span class="input-prefix">https://</span>
                <input type="url" id="longUrl" placeholder="deine-webseite.com" autofocus>
            </div>
            
            <div class="options">
                <div class="option-group">
                    <label>Eigener Code (optional)</label>
                    <input type="text" id="customCode" placeholder="z.B. mein-link, sommer-sale-2024" title="Mindestens 3 Zeichen. Erlaubt: Buchstaben, Zahlen, Bindestrich (-) und Unterstrich (_)">
                </div>
                <div class="option-group">
                    <label>Link läuft ab nach</label>
                    <select id="expiresIn" style="width: 100%; padding: 12px 14px; border: 2px solid var(--input-border); border-radius: 10px; font-size: 0.95rem; outline: none; background: var(--card-bg); color: var(--text-primary); cursor: pointer; transition: border-color 0.2s;">
                        <option value="1">1 Stunde</option>
                        <option value="24" selected>24 Stunden (1 Tag)</option>
                        <option value="72">3 Tage</option>
                        <option value="168">7 Tage</option>
                    </select>
                </div>
            </div>
            
            <button class="btn-primary" onclick="shortenUrl()">
                ✨ Kostenlosen Kurzlink erstellen
            </button>
            
            <div id="result"></div>
        </div>
    </div>

    <footer>
        <div class="footer-content">
            <div class="version-badge">
                <span class="label">Version</span>
                <span class="version">v1.3.1</span>
            </div>
            
            <div class="footer-links">
                <a href="https://github.com/LukasLow/echter.link/blob/main/CHANGELOG.md" target="_blank">📋 Changelog</a>
                <a href="#" onclick="showAbout()">ℹ️ About</a>
                <a href="https://github.com/LukasLow/echter.link" target="_blank">🐙 GitHub</a>
            </div>
            
            <p class="footer-tech">
                Built with Go + Gin + SQLite | © 2026 echter.link
            </p>
        </div>
    </footer>

    <script>
        // Theme handling
        const themeButtons = document.querySelectorAll('.theme-btn');
        const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');

        function getSystemTheme() {
            return prefersDarkScheme.matches ? 'dark' : 'light';
        }

        function applyTheme(theme) {
            const effectiveTheme = theme === 'auto' ? getSystemTheme() : theme;
            document.documentElement.setAttribute('data-theme', effectiveTheme);
        }

        function setActiveButton(theme) {
            themeButtons.forEach(btn => {
                btn.classList.toggle('active', btn.dataset.theme === theme);
            });
        }

        // Initialize theme
        const savedTheme = localStorage.getItem('theme') || 'auto';
        applyTheme(savedTheme);
        setActiveButton(savedTheme);

        // Theme button clicks
        themeButtons.forEach(btn => {
            btn.addEventListener('click', () => {
                const theme = btn.dataset.theme;
                localStorage.setItem('theme', theme);
                applyTheme(theme);
                setActiveButton(theme);
            });
        });

        // Listen for system theme changes
        prefersDarkScheme.addEventListener('change', (e) => {
            const currentTheme = localStorage.getItem('theme') || 'auto';
            if (currentTheme === 'auto') {
                applyTheme('auto');
            }
        });

        // URL shortening
        async function shortenUrl() {
            const url = document.getElementById('longUrl').value.trim();
            const customCode = document.getElementById('customCode').value.trim();
            const expiresIn = parseInt(document.getElementById('expiresIn').value) || 0;
            
            if (!url) {
                showResult('Bitte gib eine URL ein', true);
                return;
            }
            
            const fullUrl = url.startsWith('http://') || url.startsWith('https://') ? url : 'https://' + url;
            
            try {
                const response = await fetch('/api/shorten', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        original_url: fullUrl,
                        custom_code: customCode,
                        expires_in: expiresIn
                    })
                });
                
                const data = await response.json();
                
                if (data.short_code) {
                    let html = '<strong>Dein Kurzlink:</strong><br>';
                    html += '<a href="' + data.short_url + '" target="_blank">' + data.short_url + '</a>';
                    if (data.expires_at) {
                        html += '<br><small style="color:var(--text-secondary)">Läuft ab: ' + new Date(data.expires_at).toLocaleString('de-DE') + '</small>';
                    }
                    html += '<br><button class="btn-copy" onclick="copyToClipboard(\'' + data.short_url + '\')">📋 Kopieren</button>';
                    showResult(html, false);
                    copyToClipboard(data.short_url);
                } else {
                    showResult(data.error || 'Ein Fehler ist aufgetreten', true);
                }
            } catch (err) {
                showResult('Verbindungsfehler: ' + err.message, true);
            }
        }
        
        function showResult(content, isError) {
            const result = document.getElementById('result');
            result.innerHTML = '<div class="result ' + (isError ? 'error' : '') + '">' + content + '</div>';
        }
        
        function copyToClipboard(text) {
            navigator.clipboard.writeText(text).then(() => {
                console.log('Copied:', text);
            }).catch(err => {
                console.error('Copy failed:', err);
            });
        }
        
        function showAbout() {
            const theme = document.documentElement.getAttribute('data-theme');
            alert('echter.link v1.3\n\nEin moderner URL-Shortener mit:\n• SQLite Datenbank\n• Custom Short Codes\n• Ablaufdaten (max. 7 Tage)\n• Dark Mode (aktuell: ' + theme + ')\n• Go + Gin Backend\n\n© 2026 echter.link');
        }
        
        // Enter key support
        document.getElementById('longUrl').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') shortenUrl();
        });
    </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
