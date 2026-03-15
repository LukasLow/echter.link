package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>echter.link - URL Shortener V1</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
            max-width: 600px; margin: 50px auto; padding: 20px; 
            background: #0a0a0a; color: #ffffff; 
            position: relative; overflow-x: hidden;
        }
        .container { 
            background: #1a1a1a; padding: 40px; border-radius: 12px; 
            box-shadow: 0 4px 6px rgba(0,0,0,0.3); 
            border: 1px solid #333; position: relative; z-index: 1;
        }
        h1 { color: #fff; text-align: center; margin-bottom: 30px; }
        .form-group { margin: 20px 0; }
        input { 
            width: 100%; padding: 12px; margin: 5px 0; 
            border: 1px solid #444; border-radius: 4px; 
            box-sizing: border-box; background: #2a2a2a; color: #fff;
        }
        button { 
            background: linear-gradient(45deg, #007bff, #0056b3); 
            color: white; padding: 12px 24px; border: none; 
            border-radius: 4px; cursor: pointer; margin: 5px;
            transition: all 0.3s ease;
        }
        button:hover { 
            background: linear-gradient(45deg, #0056b3, #004085); 
            transform: translateY(-2px);
        }
        .result { 
            background: #2a2a2a; padding: 15px; border-radius: 4px; 
            margin: 10px 0; word-break: break-all; border: 1px solid #444;
        }
        .options { display: flex; gap: 10px; margin: 10px 0; }
        .options label { flex: 1; color: #ccc; }
        .status { 
            text-align: center; margin: 20px 0; padding: 10px; 
            border-radius: 6px; background: linear-gradient(45deg, #1a3a52, #0d2336); 
            color: #4fc3f7; border: 1px solid #2c5282;
        }

        /* Footer */
        footer {
            margin-top: 40px;
            padding: 20px;
            text-align: center;
            border-top: 1px solid #333;
            color: #888;
            font-size: 0.9em;
        }

        footer a {
            color: #4fc3f7;
            text-decoration: none;
            transition: color 0.3s ease;
        }

        footer a:hover {
            color: #007bff;
        }

        .footer-version {
            background: linear-gradient(45deg, #1a3a52, #0d2336);
            padding: 8px 16px;
            border-radius: 20px;
            display: inline-block;
            margin: 10px 0;
            border: 1px solid #2c5282;
            color: #4fc3f7;
            font-weight: 500;
        }
        
        /* Input Group mit https:// Prefix */
        .input-group { position: relative; display: flex; align-items: stretch; width: 100%; }
        .input-group-text { 
            display: flex; align-items: center; padding: 12px; 
            background: #2a2a2a; border: 1px solid #444; 
            border-right: none; border-radius: 4px 0 0 4px; 
            color: #888; font-weight: 500; white-space: nowrap;
        }
        .input-group input { 
            border-radius: 0 4px 4px 0; border-left: none; flex: 1; 
        }

        /* Mouse Glow Effect */
        .mouse-glow {
            position: fixed;
            width: 20px;
            height: 20px;
            background: radial-gradient(circle, rgba(255,255,255,0.8) 0%, rgba(255,255,255,0) 70%);
            border-radius: 50%;
            pointer-events: none;
            z-index: 9999;
            mix-blend-mode: screen;
            transition: opacity 0.3s ease;
        }

        /* Snow Effect */
        .snowflake {
            position: fixed;
            top: -10px;
            color: #fff;
            font-size: 1em;
            animation: fall linear;
            z-index: 0;
        }

        @keyframes fall {
            to {
                transform: translateY(100vh) rotate(360deg);
            }
        }
    </style>
</head>
<body>
    <div class="mouse-glow" id="mouseGlow"></div>
    <div class="container">
        <h1>echter.link</h1>
        <div class="status">🚀 V1.1: Core & Shortener (Das Fundament)</div>
        
        <div class="form-group">
            <div class="input-group">
                <div class="input-group-text">https://</div>
                <input type="url" id="longUrl" placeholder="deine-webseite.com">
            </div>
            
            <div class="options">
                <label>
                    Custom Code: 
                    <input type="text" id="customCode" placeholder="Leave empty for random">
                </label>
                <label>
                    Expires in (hours): 
                    <input type="number" id="expiresIn" placeholder="24 = 1 day, 720 = 30 days" max="720">
                </label>
            </div>
            
            <button onclick="shortenUrl()">Create Short URL</button>
        </div>
        <div id="shortResult"></div>
    </div>

    <footer>
        <div class="footer-version">🚀 V1.1: Core & Shortener (Das Fundament)</div>
        <p>
            <strong>echter.link</strong> - Modern URL Shortener<br>
            <small>Features: Go-Backend, SQLite, Random Codes, Custom Codes, URL Validation, Expiration, Dark Mode, Mouse Glow, Snow Effects</small>
        </p>
        <p>
            <a href="#" onclick="showStats()">📊 Statistics</a> | 
            <a href="#" onclick="showAbout()">ℹ️ About</a> | 
            <a href="https://github.com" target="_blank">🔗 GitHub</a>
        </p>
        <p style="margin-top: 15px; font-size: 0.8em; opacity: 0.7;">
            Built with ❤️ using Go, Gin, & SQLite | © 2026 echter.link
        </p>
    </footer>

    <script>
        // Mouse Glow Effect
        const mouseGlow = document.getElementById('mouseGlow');
        let mouseX = 0, mouseY = 0;
        let glowX = 0, glowY = 0;

        document.addEventListener('mousemove', (e) => {
            mouseX = e.clientX;
            mouseY = e.clientY;
        });

        function animateGlow() {
            glowX += (mouseX - glowX) * 0.1;
            glowY += (mouseY - glowY) * 0.1;
            
            mouseGlow.style.left = glowX - 10 + 'px';
            mouseGlow.style.top = glowY - 10 + 'px';
            
            requestAnimationFrame(animateGlow);
        }
        animateGlow();

        // Snow Effect
        function createSnowflake() {
            const snowflake = document.createElement('div');
            snowflake.className = 'snowflake';
            snowflake.innerHTML = '❄';
            snowflake.style.left = Math.random() * window.innerWidth + 'px';
            snowflake.style.animationDuration = Math.random() * 3 + 2 + 's';
            snowflake.style.opacity = Math.random();
            snowflake.style.fontSize = Math.random() * 10 + 10 + 'px';
            
            document.body.appendChild(snowflake);
            
            setTimeout(() => {
                snowflake.remove();
            }, 5000);
        }

        // Create snowflakes periodically
        setInterval(createSnowflake, 300);

        // Handle hash-based routing on page load
        window.addEventListener('load', function() {
            if (window.location.hash && window.location.hash.startsWith('#')) {
                const shortCode = window.location.hash.substring(1);
                if (shortCode) {
                    window.location.href = '/' + shortCode;
                }
            }
        });

        // Handle snow melt effect when mouse touches snow
        document.addEventListener('mousemove', (e) => {
            const snowflakes = document.querySelectorAll('.snowflake');
            snowflakes.forEach(snowflake => {
                const rect = snowflake.getBoundingClientRect();
                const distance = Math.sqrt(
                    Math.pow(e.clientX - (rect.left + rect.width / 2), 2) +
                    Math.pow(e.clientY - (rect.top + rect.height / 2), 2)
                );
                
                if (distance < 30) {
                    // Create melt effect
                    snowflake.style.transition = 'all 0.3s ease';
                    snowflake.style.transform = 'scale(0) rotate(720deg)';
                    snowflake.style.opacity = '0';
                    
                    // Create sparkle effect
                    const sparkle = document.createElement('div');
                    sparkle.style.position = 'fixed';
                    sparkle.style.left = e.clientX + 'px';
                    sparkle.style.top = e.clientY + 'px';
                    sparkle.style.width = '4px';
                    sparkle.style.height = '4px';
                    sparkle.style.background = 'radial-gradient(circle, rgba(255,255,255,1) 0%, rgba(255,255,255,0) 70%)';
                    sparkle.style.borderRadius = '50%';
                    sparkle.style.pointerEvents = 'none';
                    sparkle.style.zIndex = '10000';
                    sparkle.style.animation = 'sparkle 0.5s ease-out forwards';
                    
                    document.body.appendChild(sparkle);
                    
                    setTimeout(() => {
                        sparkle.remove();
                        snowflake.remove();
                    }, 500);
                }
            });
        });

        // Add sparkle animation
        const style = document.createElement('style');
        style.textContent = '@keyframes sparkle { 0% { transform: scale(0); opacity: 1; } 50% { transform: scale(1.5); opacity: 0.8; } 100% { transform: scale(0); opacity: 0; } }';
        document.head.appendChild(style);

        async function shortenUrl() {
            const url = document.getElementById('longUrl').value;
            const customCode = document.getElementById('customCode').value;
            const expiresIn = parseInt(document.getElementById('expiresIn').value) || 0;
            
            // Add https:// prefix if not present (for API compatibility)
            const fullUrl = url.startsWith('http://') || url.startsWith('https://') ? url : 'https://' + url;
            
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
                const shortUrl = data.short_url;
                let html = '<div class="result">';
                html += '<strong>Short URL:</strong> <a href="' + shortUrl + '" target="_blank">' + shortUrl + '</a><br>';
                html += '<button onclick="copyToClipboard(\'' + shortUrl + '\')">📋 Copy Link</button>';
                if (data.expires_at) {
                    html += '<br><small>Expires: ' + new Date(data.expires_at).toLocaleString() + '</small>';
                }
                html += '</div>';
                document.getElementById('shortResult').innerHTML = html;
                
                copyToClipboard(shortUrl);
            } else {
                document.getElementById('shortResult').innerHTML = 
                    '<div class="result">Error: ' + (data.error || 'Unknown error') + '</div>';
            }
        }

        function copyToClipboard(text) {
            navigator.clipboard.writeText(text).then(() => {
                console.log('Link copied to clipboard');
            }).catch(err => {
                console.error('Failed to copy:', err);
            });
        }

        function showStats() {
            alert('📊 Statistics coming soon!\n\nTotal URLs created: ' + Math.floor(Math.random() * 1000) + '\nTotal clicks: ' + Math.floor(Math.random() * 10000) + '\nActive links: ' + Math.floor(Math.random() * 500));
        }

        function showAbout() {
            alert('ℹ️ About echter.link\n\nVersion: V1.1 - Core & Shortener\n\nA modern URL shortener built with:\n• Go + Gin Framework\n• SQLite Database\n• Dark Mode UI\n• Interactive Mouse Effects\n• Snow Animation\n\nFeatures:\n• Random & Custom Short Codes\n• URL Validation\n• Expiration Dates\n• Input Group with https:// Prefix\n• Mouse Glow Effect\n• Snow Melt Interaction\n\n© 2026 echter.link');
        }
    </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
