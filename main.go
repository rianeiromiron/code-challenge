package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type GeocodingResponse struct {
	Results []Coordenadas `json:"results"`
}

type Coordenadas struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Clima struct {
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Timeseries []Timeserie `json:"timeseries"`
	} `json:"properties"`
}

type Timeserie struct {
	Time string `json:"time"`
	Data struct {
		Instant struct {
			Details struct {
				AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
				AirTemperature        float64 `json:"air_temperature"`
				CloudAreaFraction     float64 `json:"cloud_area_fraction"`
				RelativeHumidity      float64 `json:"relative_humidity"`
				WindFromDirection     float64 `json:"wind_from_direction"`
				WindSpeed             float64 `json:"wind_speed"`
			} `json:"details"`
		} `json:"instant"`
		Next12Hours struct {
			Summary struct {
				SymbolCode string `json:"symbol_code"`
			} `json:"summary"`
			Details struct {
				PrecipitationAmount float64 `json:"precipitation_amount"`
			} `json:"details"`
		} `json:"next_12_hours"`
		Next1Hours struct {
			Summary struct {
				SymbolCode string `json:"symbol_code"`
			} `json:"summary"`
			Details struct {
				PrecipitationAmount float64 `json:"precipitation_amount"`
			} `json:"details"`
		} `json:"next_1_hours"`
		Next6Hours struct {
			Summary struct {
				SymbolCode string `json:"symbol_code"`
			} `json:"summary"`
			Details struct {
				PrecipitationAmount float64 `json:"precipitation_amount"`
			} `json:"details"`
		} `json:"next_6_hours"`
	} `json:"data"`
}

type ClimaResponse struct {
	Pais       string      `json:"pais"`
	Latitud    float64     `json:"latitud"`
	Longitud   float64     `json:"longitud"`
	Resultados []Resultado `json:"resultados"`
}

type Resultado struct {
	FechaHora        string  `json:"fecha_hora"`
	Temperatura      float64 `json:"temperatura_c"`
	Presion          float64 `json:"presion_hpa"`
	Humedad          float64 `json:"humedad_pct"`
	Viento           float64 `json:"viento_ms"`
	VientoDir        float64 `json:"viento_dir"`
	Nubes            float64 `json:"nubes_pct"`
	Simbolo12h       string  `json:"simbolo_12h,omitempty"`
	Precipitacion12h float64 `json:"precipitacion_12h,omitempty"`
	Simbolo1h        string  `json:"simbolo_1h,omitempty"`
	Precipitacion1h  float64 `json:"precipitacion_1h,omitempty"`
	Simbolo6h        string  `json:"simbolo_6h,omitempty"`
	Precipitacion6h  float64 `json:"precipitacion_6h,omitempty"`
}

func obtenerClima(ctx context.Context, pais string, limit int) (*ClimaResponse, error) {
	query := url.Values{}
	query.Set("name", pais)
	query.Set("language", "es")

	apiURL := "https://geocoding-api.open-meteo.com/v1/search?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error geocoding: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var localizacion GeocodingResponse
	if err := json.Unmarshal(body, &localizacion); err != nil {
		return nil, fmt.Errorf("error parseando geocoding: %w", err)
	}
	if len(localizacion.Results) == 0 {
		return nil, fmt.Errorf("país no encontrado: %s", pais)
	}

	lat := localizacion.Results[0].Latitude
	lon := localizacion.Results[0].Longitude

	query = url.Values{}
	query.Set("lat", fmt.Sprint(lat))
	query.Set("lon", fmt.Sprint(lon))

	apiURL = "https://api.met.no/weatherapi/locationforecast/2.0/compact.json?" + query.Encode()
	req, err = http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "PruebasClima/1.0 rianeiromiron@gmail.com")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error API clima: %w", err)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var datosclima Clima
	if err := json.Unmarshal(body, &datosclima); err != nil {
		return nil, fmt.Errorf("error parseando clima: %w", err)
	}

	if limit > len(datosclima.Properties.Timeseries) || limit <= 0 {
		limit = len(datosclima.Properties.Timeseries)
	}

	resp := &ClimaResponse{
		Pais:     pais,
		Latitud:  lat,
		Longitud: lon,
	}

	for _, t := range datosclima.Properties.Timeseries[:limit] {
		r := Resultado{
			FechaHora:   t.Time,
			Temperatura: t.Data.Instant.Details.AirTemperature,
			Presion:     t.Data.Instant.Details.AirPressureAtSeaLevel,
			Humedad:     t.Data.Instant.Details.RelativeHumidity,
			Viento:      t.Data.Instant.Details.WindSpeed,
			VientoDir:   t.Data.Instant.Details.WindFromDirection,
			Nubes:       t.Data.Instant.Details.CloudAreaFraction,
		}
		if t.Data.Next12Hours.Summary.SymbolCode != "" {
			r.Simbolo12h = t.Data.Next12Hours.Summary.SymbolCode
			r.Precipitacion12h = t.Data.Next12Hours.Details.PrecipitationAmount
		}
		if t.Data.Next1Hours.Summary.SymbolCode != "" {
			r.Simbolo1h = t.Data.Next1Hours.Summary.SymbolCode
			r.Precipitacion1h = t.Data.Next1Hours.Details.PrecipitationAmount
		}
		if t.Data.Next6Hours.Summary.SymbolCode != "" {
			r.Simbolo6h = t.Data.Next6Hours.Summary.SymbolCode
			r.Precipitacion6h = t.Data.Next6Hours.Details.PrecipitationAmount
		}
		resp.Resultados = append(resp.Resultados, r)
	}

	return resp, nil
}

func runCLI(pais string, limit int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resultado, err := obtenerClima(ctx, pais, limit)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("\n=== %s (%.4f, %.4f) ===\n", resultado.Pais, resultado.Latitud, resultado.Longitud)
	for i, r := range resultado.Resultados {
		fmt.Printf("\n--- #%d ---\n", i+1)
		fmt.Printf("Hora: %s\n", r.FechaHora)
		fmt.Printf("Temp: %.1f°C | Presión: %.1f hPa | Humedad: %.1f%%\n", r.Temperatura, r.Presion, r.Humedad)
		fmt.Printf("Viento: %.1f m/s (%.0f°) | Nubes: %.1f%%\n", r.Viento, r.VientoDir, r.Nubes)
	}
}

func runCLIInteractive() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Ingrese el país a buscar: ")
	pais, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error leyendo input:", err)
		return
	}
	pais = strings.TrimSpace(pais)

	fmt.Print("Ingrese el número de resultados a buscar: ")
	limitStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error leyendo input:", err)
		return
	}
	limitStr = strings.TrimSpace(limitStr)
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		fmt.Println("Error: ingrese un número válido")
		return
	}

	runCLI(pais, limit)
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

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("index").Parse(indexTemplate))
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

func climaHandler(w http.ResponseWriter, r *http.Request) {
	pais := r.URL.Query().Get("pais")
	if pais == "" {
		http.Error(w, `{"error":"parámetro 'pais' requerido"}`, http.StatusBadRequest)
		return
	}

	limit := 1
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	resultado, err := obtenerClima(ctx, pais, limit)
	if err != nil {
		code := http.StatusBadGateway
		if strings.Contains(err.Error(), "país no encontrado") {
			code = http.StatusNotFound
		}
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func runServer(addr string) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/clima", enableCORS(climaHandler))

	fmt.Printf("🌤️  Servidor iniciado en http://localhost%s\n", addr)
	fmt.Printf("   Abre tu navegador en http://localhost%s\n", addr)
	fmt.Printf("   API: curl \"http://localhost%s/clima?pais=Mexico&limit=2\"\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	var serverAddr string
	var pais string
	var limit int

	flag.StringVar(&serverAddr, "server", "", "Iniciar modo servidor REST (ej: :8080)")
	flag.StringVar(&pais, "pais", "", "País a buscar")
	flag.IntVar(&limit, "limit", 1, "Número de resultados")
	flag.Parse()

	if serverAddr != "" {
		runServer(serverAddr)
	} else if pais != "" {
		runCLI(pais, limit)
	} else {
		runCLIInteractive()
	}
}