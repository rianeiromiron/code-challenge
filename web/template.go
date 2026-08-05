package web

import "html/template"

var Tmpl *template.Template

func init() {
	Tmpl = template.Must(template.New("index").Parse(indexTemplate))
}

const indexTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Clima API</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 700px;
            margin: 0 auto;
        }
        .header {
            text-align: center;
            color: white;
            margin-bottom: 30px;
            padding-top: 20px;
        }
        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 8px;
            text-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }
        .header p {
            opacity: 0.9;
            font-size: 1.1rem;
        }
        .search-card {
            background: white;
            border-radius: 16px;
            padding: 24px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            margin-bottom: 24px;
        }
        .search-row {
            display: flex;
            gap: 12px;
            flex-wrap: wrap;
        }
        .search-row input {
            flex: 1;
            min-width: 200px;
            padding: 14px 18px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 16px;
            transition: border-color 0.2s;
        }
        .search-row input:focus {
            outline: none;
            border-color: #667eea;
        }
        .search-row select {
            padding: 14px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 16px;
            background: white;
            cursor: pointer;
        }
        .search-row button {
            padding: 14px 28px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 10px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.1s, box-shadow 0.2s;
        }
        .search-row button:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }
        .search-row button:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }
        .coords {
            text-align: center;
            margin: 16px 0 8px;
        }
        .coords span {
            display: inline-block;
            background: #e3f2fd;
            color: #1565c0;
            padding: 6px 16px;
            border-radius: 20px;
            font-size: 14px;
            font-weight: 500;
        }
        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        .spinner {
            width: 40px;
            height: 40px;
            border: 4px solid #e0e0e0;
            border-top-color: #667eea;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
            margin: 0 auto 16px;
        }
        @keyframes spin { to { transform: rotate(360deg); } }
        .error-box {
            background: #ffebee;
            color: #c62828;
            padding: 16px;
            border-radius: 10px;
            text-align: center;
            font-weight: 500;
        }
        .results {
            display: grid;
            gap: 16px;
        }
        .weather-card {
            background: white;
            border-radius: 14px;
            padding: 22px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.08);
            border-left: 5px solid #667eea;
            animation: slideIn 0.3s ease-out;
        }
        @keyframes slideIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 14px;
            flex-wrap: wrap;
            gap: 8px;
        }
        .card-header h3 {
            font-size: 15px;
            color: #555;
            font-weight: 500;
        }
        .temp {
            font-size: 28px;
            font-weight: 800;
            color: #e65100;
        }
        .metrics {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 10px;
        }
        .metric {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 10px 14px;
            background: #f8f9fa;
            border-radius: 10px;
        }
        .metric-icon {
            font-size: 20px;
        }
        .metric-label {
            font-size: 11px;
            color: #888;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .metric-value {
            font-size: 15px;
            font-weight: 700;
            color: #333;
        }
        .precip {
            margin-top: 14px;
            padding-top: 14px;
            border-top: 1px solid #eee;
            display: flex;
            gap: 16px;
            flex-wrap: wrap;
            font-size: 13px;
            color: #666;
        }
        .precip strong {
            color: #1565c0;
        }
        .footer {
            text-align: center;
            color: rgba(255,255,255,0.7);
            margin-top: 30px;
            font-size: 13px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🌤️ Clima API</h1>
            <p>Consulta el pronóstico meteorológico de cualquier país</p>
        </div>

        <div class="search-card">
            <div class="search-row">
                <input type="text" id="pais" placeholder="Ej: Mexico, España, Argentina..." value="Mexico">
                <select id="limit">
                    <option value="1">1 resultado</option>
                    <option value="3" selected>3 resultados</option>
                    <option value="5">5 resultados</option>
                    <option value="10">10 resultados</option>
                </select>
                <button id="btn" onclick="buscar()">Buscar</button>
            </div>
        </div>

        <div id="resultado"></div>

        <div class="footer">
            Datos proporcionados por Open-Meteo & MET Norway
        </div>
    </div>

    <script>
        async function buscar() {
            const pais = document.getElementById('pais').value.trim();
            const limit = document.getElementById('limit').value;
            const btn = document.getElementById('btn');
            const out = document.getElementById('resultado');

            if (!pais) {
                out.innerHTML = '<div class="error-box">Ingresa un país</div>';
                return;
            }

            btn.disabled = true;
            out.innerHTML = '<div class="loading"><div class="spinner"></div><div>Consultando pronóstico...</div></div>';

            try {
                const res = await fetch('/clima?pais=' + encodeURIComponent(pais) + '&limit=' + limit);
                const data = await res.json();

                if (!res.ok) {
                    throw new Error(data.error || 'Error del servidor');
                }

                render(data);
            } catch (e) {
                out.innerHTML = '<div class="error-box">' + e.message + '</div>';
            } finally {
                btn.disabled = false;
            }
        }

        function render(data) {
            const out = document.getElementById('resultado');
            let html = '<div class="coords"><span>📍 ' + data.latitud.toFixed(4) + ', ' + data.longitud.toFixed(4) + '</span></div>';
            html += '<div class="results">';

            data.resultados.forEach(r => {
                const fecha = new Date(r.fecha_hora).toLocaleString('es-ES', {
                    weekday: 'short', day: 'numeric', month: 'short',
                    hour: '2-digit', minute: '2-digit'
                });

                html += '<div class="weather-card">';
                html += '<div class="card-header"><h3>🕐 ' + fecha + '</h3><span class="temp">' + r.temperatura_c.toFixed(1) + '°C</span></div>';
                html += '<div class="metrics">';
                html += '<div class="metric"><span class="metric-icon">💨</span><div><div class="metric-label">Viento</div><div class="metric-value">' + r.viento_ms.toFixed(1) + ' m/s</div></div></div>';
                html += '<div class="metric"><span class="metric-icon">💧</span><div><div class="metric-label">Humedad</div><div class="metric-value">' + r.humedad_pct.toFixed(0) + '%</div></div></div>';
                html += '<div class="metric"><span class="metric-icon">📊</span><div><div class="metric-label">Presión</div><div class="metric-value">' + r.presion_hpa.toFixed(1) + ' hPa</div></div></div>';
                html += '<div class="metric"><span class="metric-icon">☁️</span><div><div class="metric-label">Nubes</div><div class="metric-value">' + r.nubes_pct.toFixed(0) + '%</div></div></div>';
                html += '</div>';
                html += '<div class="precip">';
                if (r.simbolo_1h) html += '<div>1h: <strong>' + r.simbolo_1h + '</strong> (' + r.precipitacion_1h + ' mm)</div>';
                if (r.simbolo_6h) html += '<div>6h: <strong>' + r.simbolo_6h + '</strong> (' + r.precipitacion_6h + ' mm)</div>';
                if (r.simbolo_12h) html += '<div>12h: <strong>' + r.simbolo_12h + '</strong> (' + r.precipitacion_12h + ' mm)</div>';
                html += '</div></div>';
            });

            html += '</div>';
            out.innerHTML = html;
        }

        document.getElementById('pais').addEventListener('keypress', e => {
            if (e.key === 'Enter') buscar();
        });
    </script>
</body>
</html>`
