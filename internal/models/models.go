package models

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
