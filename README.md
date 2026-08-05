# 🌤️ Clima API

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Aplicación dual escrita en **Go** que consulta el pronóstico meteorológico de cualquier país mediante APIs externas. Puede ejecutarse como **CLI interactivo**, **CLI con flags** o como **servidor web REST** con interfaz HTML embebida.

> Datos proporcionados por [Open-Meteo](https://open-meteo.com/) (geocoding) y [MET Norway](https://api.met.no/) (forecast).

---

## ✨ Características

- 🖥️ **Tres modos de ejecución**: CLI interactivo, CLI con argumentos y servidor HTTP.
- 🌐 **API REST** con endpoint `/clima` que devuelve JSON estructurado.
- 🎨 **Frontend HTML embebido** con diseño responsive (sin dependencias externas).
- 🌍 **Geocoding automático**: busca coordenadas por nombre de país.
- ⏱️ **Timeouts configurables** en todas las llamadas a APIs externas.
- 🗂️ **Arquitectura modular** siguiendo estándares de la comunidad Go.

---

## 📁 Estructura del Proyecto

code-challenge/
├── cmd/
│   ├── cli/              # Punto de entrada para modo consola
│   └── server/           # Punto de entrada para modo servidor web
├── internal/
│   ├── api/              # Lógica de negocio y llamadas a APIs externas
│   ├── handlers/         # HTTP handlers (CORS, rutas, JSON)
│   └── models/           # Structs y tipos de datos
├── web/
│   └── template.go       # Template HTML embebido
├── go.mod
├── .gitignore
└── README.md


---

## 🚀 Requisitos

- [Go](https://golang.org/dl/) **1.24** o superior
- Conexión a Internet (para consultar APIs externas)

---

## ⚙️ Instalación

```bash
# Clonar el repositorio
git clone https://github.com/rianeiromiron/code-challenge.git
cd code-challenge

# Descargar dependencias
go mod tidy

# Compilar ambos ejecutables
go build -o code-challenge-cli.exe ./cmd/cli
go build -o code-challenge-server.exe ./cmd/server

🖥️ Uso
Modo CLI Interactivo

    .\code-challenge-cli.exe

El programa te pedirá el país y el número de resultados por consola.

Modo CLI con Flags
    .\code-challenge-cli.exe -pais Mexico -limit 3
    .\code-challenge-cli.exe -pais "Costa Rica" -limit 5

Modo Servidor Web
    .\code-challenge-server.exe

Abre tu navegador en: http://localhost:8080

Endpoints
    | Método | Ruta                         | Descripción                |
    | ------ | ---------------------------- | -------------------------- |
    | `GET`  | `/`                          | Interfaz web HTML          |
    | `GET`  | `/clima?pais=Mexico&limit=3` | API REST que devuelve JSON |

Ejemplo con cURL

    curl "http://localhost:8080/clima?pais=Mexico&limit=2"

Respuesta JSON

    {
    "pais": "Mexico",
    "latitud": 19.4326,
    "longitud": -99.1332,
    "resultados": [
        {
        "fecha_hora": "2026-08-04T12:00:00Z",
        "temperatura_c": 24.5,
        "presion_hpa": 1013.2,
        "humedad_pct": 65.0,
        "viento_ms": 3.2,
        "viento_dir": 180.0,
        "nubes_pct": 40.0,
        "simbolo_1h": "clearsky_day",
        "precipitacion_1h": 0.0,
        "simbolo_6h": "partlycloudy_day",
        "precipitacion_6h": 0.3,
        "simbolo_12h": "partlycloudy_day",
        "precipitacion_12h": 0.5
        }
    ]
    }

🛠️ Tecnologías
Lenguaje: Go 1.24
Librería estándar: net/http, html/template, encoding/json, context
APIs externas: Open-Meteo Geocoding API, MET Norway Locationforecast
Frontend: HTML5 + CSS3 + JavaScript vanilla (embebido)

📌 Notas
El servidor incluye headers CORS habilitados para consumo desde cualquier origen.
Todas las peticiones a APIs externas tienen un timeout de 15 segundos.
Se requiere un User-Agent válido para la API de MET Norway (incluido en el código).

🤝 Contribuciones
Las contribuciones son bienvenidas. Si encuentras un bug o tienes una idea de mejora, abre un Issue o envía un Pull Request.

📝 Licencia
Este proyecto está bajo la licencia MIT.

👤 Autor
Rianeiro Miron — @rianeiromiron